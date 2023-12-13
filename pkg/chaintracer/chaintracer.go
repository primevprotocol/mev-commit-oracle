package chaintracer

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"

	"io"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/primevprotocol/oracle/pkg/rollupclient"
	"github.com/rs/zerolog/log"
)

type TransactionHash []byte

type BlockDetails struct {
	Transactions []string `json:"transactions"`
	BlockNumber  string   `json:"number"`
}

type InfuraResponse struct {
	Jsonrpc string       `json:"jsonrpc"`
	ID      int          `json:"id"`
	Result  BlockDetails `json:"result"`
}

type L1DataRetriver interface {
	GetTransactions(blockNumber int64) (*BlockDetails, error)
	GetWinningBuilder(blockNumber int64) (string, error)
}

// We can maintain a skipped block list in the smart contract
type PreConfirmationsContract interface {
	GetCommitmentsByBlockNumber(opts *bind.CallOpts, blockNumber *big.Int) ([][32]byte, error)
	GetTxnHashFromCommitment(opts *bind.CallOpts, commitmentIndex [32]byte) (string, error)
}

// CommitmentsStore is an interface that is used to retrieve commitments from the smart contract
// and store them in a local database
type CommitmentsStore interface {
	UpdateCommitmentsForBlockNumber(blockNumber int64) (done chan struct{}, err chan error)
	RetrieveCommitments(blockNumber int64) (map[string]bool, error)

	// Used for restarting the Commitment Store on startup
	LargestStoredBlockNumber() (int64, error)
}

type DBTxnFilter struct {
	db            *sql.DB
	preConfClient PreConfirmationsContract
}

func NewDBTxnFilter(db *sql.DB, preConfClient PreConfirmationsContract) CommitmentsStore {
	return &DBTxnFilter{
		db:            db,
		preConfClient: preConfClient,
	}
}

func (f DBTxnFilter) LargestStoredBlockNumber() (int64, error) {
	var largestBlockNumber int64
	err := f.db.QueryRow("SELECT MAX(block_number) FROM committed_transactions").Scan(&largestBlockNumber)
	if err != nil {
		return 0, err
	}
	return largestBlockNumber, nil
}

func (f DBTxnFilter) UpdateCommitmentsForBlockNumber(blockNumber int64) (done chan struct{}, errorC chan error) {
	done = make(chan struct{})
	errorC = make(chan error)

	go func(done chan struct{}, errorC chan error) {
		commitments, err := f.preConfClient.GetCommitmentsByBlockNumber(&bind.CallOpts{
			Pending: false,
			From:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
			Context: context.Background(),
		}, big.NewInt(blockNumber))
		if err != nil {
			log.Error().Err(err).Msg("Error getting commitments")
			errorC <- err
			return
		}
		log.Info().Int("block_number", int(blockNumber)).Int("commitments", len(commitments)).Msg("Retrieved commitments")
		for _, commitment := range commitments {
			commitmentTxnHash, err := f.preConfClient.GetTxnHashFromCommitment(&bind.CallOpts{
				Pending: false,
				From:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
				Context: context.Background(),
			}, commitment)
			if err != nil {
				log.Error().Err(err).Msg("Error getting txn hash from commitment")
				errorC <- err
				return
			}

			sqlStatement := `
			INSERT INTO committed_transactions (transaction, block_number)
			VALUES ($1, $2)`
			result, err := f.db.Exec(sqlStatement, commitmentTxnHash, blockNumber)
			if err != nil {
				if err, ok := err.(*pq.Error); ok {
					// Check if the error is a duplicate key violation
					if err.Code.Name() == "unique_violation" {
						log.Info().Msg("Duplicate key violation - ignoring")
						continue
					}
				}
				log.Error().Err(err).Msg("Error inserting txn into DB")
				errorC <- err
				return
			}
			rowsImpacted, err := result.RowsAffected()
			if err != nil {
				log.Error().Err(err).Msg("Error getting rows impacted")
				errorC <- err
				return
			}
			log.Info().Int("rows_affected", int(rowsImpacted)).Msg("Inserted txn into DB")
		}
		done <- struct{}{}
	}(done, errorC)

	return done, errorC
}

func (f DBTxnFilter) RetrieveCommitments(blockNumber int64) (map[string]bool, error) {
	filter := make(map[string]bool)

	rows, err := f.db.Query("SELECT transaction FROM committed_transactions WHERE block_number = $1", blockNumber)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var txnHash string
		err = rows.Scan(&txnHash)
		if err != nil {
			return nil, err
		}
		filter[txnHash] = true
	}

	return filter, nil
}

type InMemoryTxnFilter struct {
	db            map[int64][]string // BlockNumber -> TxnHashes
	preConfClient PreConfirmationsContract
}

func NewTransactionCommitmentFilter(preConfClient PreConfirmationsContract) CommitmentsStore {
	return &InMemoryTxnFilter{
		db:            make(map[int64][]string),
		preConfClient: preConfClient,
	}
}

// Reduant as InMemoryTxnFilter is not persisten
func (f InMemoryTxnFilter) LargestStoredBlockNumber() (int64, error) {
	var largestBlockNumber int64
	for blockNumber := range f.db {
		if blockNumber > largestBlockNumber {
			largestBlockNumber = blockNumber
		}
	}
	return largestBlockNumber, nil
}

func (f InMemoryTxnFilter) UpdateCommitmentsForBlockNumber(blockNumber int64) (done chan struct{}, errorC chan error) {
	done = make(chan struct{})
	errorC = make(chan error)
	if _, ok := f.db[blockNumber]; !ok {
		f.db[blockNumber] = []string{}
	}

	go func(done chan struct{}, errorC chan error) {
		commitments, err := f.preConfClient.GetCommitmentsByBlockNumber(&bind.CallOpts{
			Pending: false,
			From:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
			Context: context.Background(),
		}, big.NewInt(blockNumber))
		if err != nil {
			log.Error().Err(err).Msg("Error getting commitments")
			errorC <- err
			return
		}

		for _, commitment := range commitments {
			commitmentTxnHash, err := f.preConfClient.GetTxnHashFromCommitment(&bind.CallOpts{
				Pending: false,
				From:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
				Context: context.Background(),
			}, commitment)
			if err != nil {
				log.Error().Err(err).Msg("Error getting txn hash from commitment")
				errorC <- err
				return
			}

			f.db[blockNumber] = append(f.db[blockNumber], commitmentTxnHash)
			// Must clear this variable
			done <- struct{}{}
		}
	}(done, errorC)

	return done, errorC
}

// TODO(@ckartik): Adds init filter
// RetrieveCommitments initializes the a goroutine that will fetch all the commitments from the smart contract for the given blockNumber
// and return a channel that will be used to filter out transactions that have already been confirmed
// Need to model the creation of pre-confirmations from a searcher
// NOTE: Need to manage situation where the contracts receive a commitment after the block has been updated to blockNumber
func (f InMemoryTxnFilter) RetrieveCommitments(blockNumber int64) (filter map[string]bool, err error) {
	filter = make(map[string]bool)
	for _, txn := range f.db[blockNumber] {
		filter[txn] = true
	}

	return filter, nil
}

// SmartContractTracer is a tracer that uses the smart contract
// to retrieve details about the next block that needs to be proccesed
// it has the option of
type IntegerationTestTracer struct {
	contractClient           *rollupclient.Rollupclient
	currentBlockNumberCached int64
	cs                       CommitmentsStore

	RateLimit time.Duration
}

// SmartContractTracer is a tracer that uses the smart contract
// to retrieve details about the next block that needs to be proccesed
// it has the option of
type SmartContractTracer struct {
	contractClient           *rollupclient.Rollupclient
	currentBlockNumberCached int64

	RateLimit time.Duration
}

func (st *SmartContractTracer) GetNextBlockNumber(ctx context.Context) (NewBlockNumber int64) {
	nextBlockNumber, err := st.contractClient.GetNextRequestedBlockNumber(&bind.CallOpts{
		Pending: false,
		From:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
		Context: ctx,
	})
	for err != nil {
		select {
		case <-ctx.Done():
			log.Info().Msg("Context cancelled, exiting GetNextBlockNumber")
			return -1
		default:
			log.Error().Err(err).Msg("Error getting next block number, will go to sleep for 5 seconds and try again")
			time.Sleep(5 * time.Second)
			nextBlockNumber, err = st.contractClient.GetNextRequestedBlockNumber(&bind.CallOpts{
				Pending: false,
				From:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
				Context: ctx,
			})
		}
	}
	st.currentBlockNumberCached = nextBlockNumber.Int64()
	return st.currentBlockNumberCached
}

// TODO(@ckartik): Move logic for service based data request to an isolated function.
func (st *SmartContractTracer) RetrieveDetails() (block *BlockDetails, BlockBuilder string, err error) {
	retries := 0
	var blockData *BlockDetails
	log.Debug().Msg("Starting Retreival of Block Details")
	// Retrieve Block Details from Infura
	for blockData == nil {
		log.Debug().Msg("fetching infura data")
		blockData = InfuraData(st.currentBlockNumberCached, "https://mainnet.infura.io/v3/b8877b173a0543bea7dca82c313e7347")
		time.Sleep(time.Duration(retries*5) * time.Second)
		if retries > 5 {
			log.Error().Msg("Error: Could not retrieve block data")
			return nil, "", errors.New("Error: Could not retrieve block data")
		}
		retries += 1
	}
	retries = 0
	var builderName string
	for builderName == "" {
		log.Debug().Msg("fetching Payloadsde data")
		builderName, err = PayloadsDe(st.currentBlockNumberCached)
		if err != nil {
			log.Error().Err(err).Msg("Error: Could not retrieve block data")
		}
		time.Sleep(time.Duration(retries*5) * time.Second)
		if retries > 5 {
			return nil, "", errors.New("Error: Could not retrieve block data")
		}
		retries += 1
	}
	log.Info().
		Int64("BlockNumber", st.currentBlockNumberCached).
		Str("Builder", builderName).
		Int("Txns Received", len(blockData.Transactions)).
		Msg("Finished Retreival of Block Details")

	return blockData, builderName, nil
}

func NewIntegrationTestTracer(ctx context.Context, contractClient *rollupclient.Rollupclient, cs CommitmentsStore) Tracer {
	tracer := &IntegerationTestTracer{
		contractClient: contractClient,
		cs:             cs,
	}
	tracer.GetNextBlockNumber(ctx)
	return tracer
}

func (st *IntegerationTestTracer) GetNextBlockNumber(ctx context.Context) (NewBlockNumber int64) {
	nextBlockNumber, err := st.contractClient.GetNextRequestedBlockNumber(&bind.CallOpts{
		Pending: false,
		From:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
		Context: ctx,
	})
	for err != nil {
		select {
		case <-ctx.Done():
			log.Info().Msg("Context cancelled, exiting GetNextBlockNumber")
			return -1
		default:
			log.Error().Err(err).Msg("Error getting next block number, will go to sleep for 5 seconds and try again")
			time.Sleep(5 * time.Second)
			nextBlockNumber, err = st.contractClient.GetNextRequestedBlockNumber(&bind.CallOpts{
				Pending: false,
				From:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
				Context: ctx,
			})
		}
	}
	st.currentBlockNumberCached = nextBlockNumber.Int64()
	return st.currentBlockNumberCached
}

// TODO(@ckartik): Move logic for service based data request to an isolated function.
func (it *IntegerationTestTracer) RetrieveDetails() (block *BlockDetails, BlockBuilder string, err error) {
	txns, err := it.cs.RetrieveCommitments(it.currentBlockNumberCached)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving txns in integreation test tracer")
		return nil, "", err
	}

	blockData := &BlockDetails{
		BlockNumber:  strconv.FormatInt(it.currentBlockNumberCached, 10),
		Transactions: []string{},
	}

	for txn := range txns {
		blockData.Transactions = append(blockData.Transactions, txn)
	}
	return blockData, "dummy builder", nil
}

func NewSmartContractTracer(contractClient *rollupclient.Rollupclient, ctx context.Context) Tracer {
	tracer := &SmartContractTracer{
		contractClient: contractClient,
	}
	tracer.GetNextBlockNumber(ctx)
	return tracer

}
func NewIncrementingTracer(startingBlockNumber int64, rateLimit time.Duration) Tracer {
	return &IncrementingTracer{
		BlockNumber: startingBlockNumber,
		RateLimit:   rateLimit,
	}
}

type IncrementingTracer struct {
	BlockNumber int64
	RateLimit   time.Duration
}

type Tracer interface {
	GetNextBlockNumber(ctx context.Context) (NewBlockNumber int64)
	RetrieveDetails() (block *BlockDetails, BlockBuilder string, err error)
}

type Payload struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

func (it *IncrementingTracer) GetNextBlockNumber(_ context.Context) (NewBlockNumber int64) {
	it.BlockNumber += 1
	log.Info().Int64("BlockNumber", it.BlockNumber).Msg("Incremented Block Number")
	return it.BlockNumber
}

func (it *IncrementingTracer) RetrieveDetails() (block *BlockDetails, BlockBuilder string, err error) {
	log.Info().Str("RateLimit", it.RateLimit.String()).Msg("Waiting for Rate Limit")
	time.Sleep(it.RateLimit)
	retries := 0
	var blockData *BlockDetails
	log.Debug().Msg("Starting Retreival of Block Details")
	// Retrieve Block Details from Infura
	for blockData == nil {
		log.Debug().Msg("fetching infura data")
		blockData = InfuraData(it.BlockNumber, "https://mainnet.infura.io/v3/b8877b173a0543bea7dca82c313e7347")
		time.Sleep(time.Duration(retries*2) * time.Second)
		if retries > 5 {
			log.Error().Msg("Error: Could not retrieve block data")
			return nil, "", errors.New("Error: Could not retrieve block data")
		}
		retries += 1
	}
	retries = 0
	var builderName string
	for builderName == "" {
		log.Debug().Msg("fetching Payloadsde data")
		builderName, err = PayloadsDe(it.BlockNumber)
		if err != nil {
			log.Error().Err(err).Msg("Error: Could not retrieve block data")
		}
		time.Sleep(time.Duration(retries*2) * time.Second)
		if retries > 5 {
			return nil, "", errors.New("Error: Could not retrieve block data")
		}
		retries += 1
	}
	log.Info().
		Int64("BlockNumber", it.BlockNumber).
		Str("Builder", builderName).
		Int("Txns Received", len(blockData.Transactions)).
		Msg("Finished Retreival of Block Details")

	return blockData, builderName, nil
}

type MainnetInfuraAndPayloadsdeDataRetriever struct {
	url string
}

func NewMainnetInfuraAndPayloadsdeDataRetriever() L1DataRetriver {
	return &MainnetInfuraAndPayloadsdeDataRetriever{
		url: "https://mainnet.infura.io/v3/b8877b173a0543bea7dca82c313e7347",
	}
}

func (m MainnetInfuraAndPayloadsdeDataRetriever) GetTransactions(blockNumber int64) (*BlockDetails, error) {
	return InfuraData(blockNumber, m.url), nil
}

func (m MainnetInfuraAndPayloadsdeDataRetriever) GetWinningBuilder(blockNumber int64) (string, error) {
	return PayloadsDe(blockNumber)
}

func PayloadsDe(blockNumber int64) (string, error) {
	url := "https://api.payload.de/block_info?block=" + strconv.FormatInt(blockNumber, 10)
	resp, err := http.Get(url)
	log.Debug().Msg("payloadsde data recieved")
	if err != nil {
		log.Error().Err(err).Str("URL", url).Msg("Error sending request")
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading response")
		return "", err
	}

	var payload struct {
		Builder string `json:"builder"`
	}
	json.Unmarshal(data, &payload)

	return payload.Builder, nil
}

func InfuraData(blockNumber int64, url string) *BlockDetails {
	hex := strconv.FormatInt(blockNumber, 16)
	// url := "https://mainnet.infura.io/v3/b8877b173a0543bea7dca82c313e7347"
	payload := Payload{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{"0x" + hex, false},
		ID:      1,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling payload")
		return nil
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Error().Err(err).Msg("Error creating request")
		return nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("", "66f7989079d54ad6988e3b083da0723a")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Error sending request")
		return nil
	}
	defer resp.Body.Close()

	var infuraresp struct {
		Result BlockDetails `json:"result"`
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading response")
		return nil
	}
	json.Unmarshal(data, &infuraresp)
	if infuraresp.Result.BlockNumber == "" {
		log.Error().Msg("Error: Block data is empty")
		return nil
	}
	return &infuraresp.Result
}
