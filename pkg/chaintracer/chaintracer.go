package chaintracer

import (
	"bytes"
	"context"
	"encoding/json"
	"math/big"

	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/primevprotocol/oracle/pkg/repository"
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

// SmartContractTracer is a tracer that uses the smart contract
// to retrieve details about the next block that needs to be proccesed
// it has the option of
type SmartContractTracer struct {
	contractClient           *rollupclient.Rollupclient
	currentBlockNumberCached int64

	RateLimit           time.Duration
	startingBlockNumber int64
	L1Client            *ethclient.Client
}

func (st *SmartContractTracer) GetNextBlockNumber(ctx context.Context) (NewBlockNumber int64) {
	if st.currentBlockNumberCached < st.startingBlockNumber {
		st.currentBlockNumberCached = st.startingBlockNumber

		return st.currentBlockNumberCached
	}

	// nextBlockNumber, err := st.contractClient.GetNextRequestedBlockNumber(&bind.CallOpts{
	// 	Pending: false,
	// 	From:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
	// 	Context: ctx,
	// })
	// for err != nil {
	// 	select {
	// 	case <-ctx.Done():
	// 		log.Info().Msg("Context cancelled, exiting GetNextBlockNumber")
	// 		return -1
	// 	default:
	// 		log.Error().Err(err).Msg("Error getting next block number, will go to sleep for 5 seconds and try again")
	// 		time.Sleep(5 * time.Second)
	// 		nextBlockNumber, err = st.contractClient.GetNextRequestedBlockNumber(&bind.CallOpts{
	// 			Pending: false,
	// 			From:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"), // TODO(@ckartik): See how we can remove this
	// 			Context: ctx,
	// 		})
	// 	}
	// }
	st.currentBlockNumberCached = st.currentBlockNumberCached + 1

	return st.currentBlockNumberCached
}

// TODO(@ckartik): Move logic for service based data request to an isolated function.
func (st *SmartContractTracer) RetrieveDetails() (block *BlockDetails, BlockBuilder string, err error) {
	l1BlockNumber, err := st.L1Client.BlockNumber(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Error getting block number")
		return nil, "", err
	}
	// Make fast mode flexible
	if l1BlockNumber > uint64(st.currentBlockNumberCached+4) {
		log.Info().Uint64("l1_block_number", l1BlockNumber).Int64("current_oracle_block", st.currentBlockNumberCached).Msg("Fast Mode")
		time.Sleep(2 * time.Second)
	} else {
		log.Info().Uint64("l1_block_number", l1BlockNumber).Int64("current_oracle_block", st.currentBlockNumberCached).Msg("Normal Mode")
		time.Sleep(12 * time.Second)
	}

	log.Debug().Msg("Starting Retreival of Block Details")
	L1Block, err := st.L1Client.BlockByNumber(context.Background(), big.NewInt(st.currentBlockNumberCached))

	// Retrieve Block Details from Infura
	for retries := 0; err != nil; retries++ {
		log.Debug().Msg("Failed to get block from L1, will try again")
		time.Sleep(time.Duration(retries*5) * time.Second)
		if retries > 5 {
			log.Error().Msg("Error: Could not retrieve block data")
			return nil, "", errors.New("Error: Could not retrieve block data")
		}
		L1Block, err = st.L1Client.BlockByNumber(context.Background(), big.NewInt(st.currentBlockNumberCached))
	}
	blockData := &BlockDetails{
		BlockNumber:  strconv.FormatInt(st.currentBlockNumberCached, 10),
		Transactions: []string{},
	}

	txns := L1Block.Transactions()
	for _, txn := range txns {
		blockData.Transactions = append(blockData.Transactions, txn.Hash().String())
	}

	return blockData, string(L1Block.Header().Extra), nil
}

func NewSmartContractTracer(ctx context.Context, contractClient *rollupclient.Rollupclient, l1Client *ethclient.Client, startingBlockNumber int64) Tracer {
	tracer := &SmartContractTracer{
		contractClient:      contractClient,
		startingBlockNumber: startingBlockNumber,
		L1Client:            l1Client,
	}

	return tracer

}

/*
Redundant Testing Code Begins here
*/
type L1DataRetriver interface {
	GetTransactions(blockNumber int64) (*BlockDetails, error)
	GetWinningBuilder(blockNumber int64) (string, error)
	GetLatestBlockNumber() (int64, error)
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

func (st *MainnetInfuraAndPayloadsdeDataRetriever) GetLatestBlockNumber() (int64, error) {
	return 0, errors.New("not implemented")
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

func NewIntegrationTestTracer(ctx context.Context, contractClient *rollupclient.Rollupclient, cs repository.CommitmentsStore) Tracer {
	tracer := &IntegerationTestTracer{
		contractClient: contractClient,
		cs:             cs,
	}
	tracer.GetNextBlockNumber(ctx)
	return tracer
}

func (st *IntegerationTestTracer) GetNextBlockNumber(ctx context.Context) (NewBlockNumber int64) {
	// nextBlockNumber, err := st.contractClient.GetNextRequestedBlockNumber(&bind.CallOpts{
	// 	Pending: false,
	// 	From:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
	// 	Context: ctx,
	// })
	// for err != nil {
	// 	select {
	// 	case <-ctx.Done():
	// 		log.Info().Msg("Context cancelled, exiting GetNextBlockNumber")
	// 		return -1
	// 	default:
	// 		log.Error().Err(err).Msg("Error getting next block number, will go to sleep for 5 seconds and try again")
	// 		time.Sleep(5 * time.Second)
	// 		nextBlockNumber, err = st.contractClient.GetNextRequestedBlockNumber(&bind.CallOpts{
	// 			Pending: false,
	// 			From:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
	// 			Context: ctx,
	// 		})
	// 	}
	// }
	st.currentBlockNumberCached += 1
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
		time.Sleep(50 * time.Millisecond)
		blockData.Transactions = append(blockData.Transactions, txn)
	}
	return blockData, "dummy builder", nil
}

// IntegerationTestTracer is a tracer that uses the smart contract
// to retrieve details about the next block that needs to be proccesed
// it has the option of
type IntegerationTestTracer struct {
	contractClient           *rollupclient.Rollupclient
	currentBlockNumberCached int64
	cs                       repository.CommitmentsStore

	RateLimit time.Duration
}
