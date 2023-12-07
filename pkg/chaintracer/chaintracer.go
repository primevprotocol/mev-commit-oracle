package chaintracer

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/primevprotocol/oracle/pkg/preconf"
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

// We can maintain a skipped block list in the smart contract

type optimizationFilter struct {
	preConfClient *preconf.PreConfClient
}

// The Future Fitler interface is used to initialize the filter
type TransactionFilter interface {
	InitFilter(blockNumber int64) (chan map[string]bool, chan error)
}

func NewTransactionCommitmentFilter(preConfClient *preconf.PreConfClient) TransactionFilter {
	return &optimizationFilter{
		preConfClient: preConfClient,
	}
}

// TODO(@ckartik): Adds init filter
// Need to model the creation of pre-confirmations from a searcher
// NOTE: Need to manage situation where the contracts receive a commitment after the block has been updated to blockNumber
func (f optimizationFilter) InitFilter(blockNumber int64) (filter chan map[string]bool, errChannel chan error) {
	filter = make(chan map[string]bool)
	errChannel = make(chan error)
	go func(filter chan map[string]bool, errChannel chan error) {
		db := make(map[string]bool)
		commitments, err := f.preConfClient.GetCommitmentsByBlockNumber(&bind.CallOpts{
			Pending: false,
			From:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
			Context: context.Background(),
		}, big.NewInt(blockNumber))
		if err != nil {
			log.Error().Err(err).Msg("Error getting commitments")
			errChannel <- err
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
				errChannel <- err
				return
			}

			// Must clear this variable
			db[commitmentTxnHash] = true // Set encountered TxnHash to true
		}
		filter <- db
	}(filter, errChannel)

	return filter, errChannel
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
		blockData = InfuraData(st.currentBlockNumberCached)
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
		blockData = InfuraData(it.BlockNumber)
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

func InfuraData(blockNumber int64) *BlockDetails {
	hex := strconv.FormatInt(blockNumber, 16)
	url := "https://mainnet.infura.io/v3/b8877b173a0543bea7dca82c313e7347"
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

func NewDummyTracer() Tracer {
	return &dummyTracer{
		blockNumberCurrent: 0,
	}
}

type dummyTracer struct {
	blockNumberCurrent int64
}

func (d *dummyTracer) GetNextBlockNumber(_ context.Context) int64 {
	d.blockNumberCurrent += 1
	return d.blockNumberCurrent
}

func (d *dummyTracer) RetrieveDetails() (block *BlockDetails, BlockBuilder string, err error) {
	block = &BlockDetails{
		BlockNumber:  strconv.FormatInt(d.blockNumberCurrent, 10),
		Transactions: []string{},
	}

	for i := 0; i < 200; i++ {
		randomInt, err := rand.Int(rand.Reader, big.NewInt(1000))
		if err != nil {
			panic(err)
		}
		randomBytes := crypto.Keccak256(randomInt.Bytes())
		block.Transactions = append(block.Transactions, hex.EncodeToString(randomBytes))
	}

	sleepDuration, _ := rand.Int(rand.Reader, big.NewInt(12))
	time.Sleep(time.Duration(sleepDuration.Int64()) * time.Second)
	return block, "k builder", nil
}
