package chaintracer

import (
	"context"
	"math/big"

	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/primevprotocol/mev-oracle/pkg/rollupclient"
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
	contractClient           *rollupclient.OracleClient
	currentBlockNumberCached int64

	RateLimit           time.Duration
	startingBlockNumber int64
	L1Client            *ethclient.Client

	// Configurable parameters fastModeSleep and normalModeSleep fastModeSensitivity

	fastModeSleep       time.Duration
	normalModeSleep     time.Duration
	fastModeSensitivity int64

	integreationTestMode bool
}

func (st *SmartContractTracer) GetNextBlockNumber(ctx context.Context) (NewBlockNumber int64) {

	// TODO(@ckartik): Use stored block number on contract instead of incrementing (For Failure reslience)
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
				From:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"), // TODO(@ckartik): See how we can remove this
				Context: ctx,
			})
		}
	}
	st.currentBlockNumberCached = nextBlockNumber.Int64()
	if nextBlockNumber.Int64() < st.startingBlockNumber {
		log.Info().Int64("next_block_number", nextBlockNumber.Int64()).Int64("starting_block_number", st.startingBlockNumber).Msg("Next block number is less than starting block number, returning starting block number")
		st.currentBlockNumberCached = st.startingBlockNumber
	}

	return st.currentBlockNumberCached
}

var IntegrationTestBuilders = []string{
	"0x48ddC642514370bdaFAd81C91e23759B0302C915",
	"0x972eb4Fc3c457da4C957306bE7Fa1976BB8F39A6",
	"0xA1e8FDB3bb6A0DB7aA5Db49a3512B01671686DCB",
	"0xB9286CB4782E43A202BfD426AbB72c8cb34f886c",
	"0xdaa1EEe546fc3f2d10C348d7fEfACE727C1dfa5B",
	"0x93DC0b6A7F454Dd10373f1BdA7Fe80BB549EE2F9",
	"0x426184Df456375BFfE7f53FdaF5cB48DeB3bbBE9",
	"0x41cC09BD5a97F22045fe433f1AF0B07d0AB28F58",
}

// TODO(@ckartik): Move logic for service based data request to an isolated function.
func (st *SmartContractTracer) RetrieveDetails() (block *BlockDetails, BlockBuilder string, err error) {
	l1BlockNumber, err := st.L1Client.BlockNumber(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Error getting block number")
		return nil, "", err
	}
	// Make fast mode flexible
	// no hardcoded variables
	if l1BlockNumber > uint64(st.currentBlockNumberCached+st.fastModeSensitivity) {
		log.Info().Uint64("l1_block_number", l1BlockNumber).Int64("current_oracle_block", st.currentBlockNumberCached).Msg("Fast Mode")
		time.Sleep(st.fastModeSleep * time.Second)
	} else {
		log.Info().Uint64("l1_block_number", l1BlockNumber).Int64("current_oracle_block", st.currentBlockNumberCached).Msg("Normal Mode")
		time.Sleep(st.normalModeSleep * time.Second)
	}

	log.Debug().Msg("Starting Retreival of Block Details")
	L1Block, err := st.L1Client.BlockByNumber(context.Background(), big.NewInt(st.currentBlockNumberCached))

	// Retrieve Block Details from Infura
	for retries := 0; err != nil; retries++ {
		// The err is being printed here.
		log.Debug().Msg("Failed to get block from L1, will try again")
		time.Sleep(time.Duration(retries*5) * time.Second)
		if retries > 5 {
			// Logging and returning is a bad pattern. Log in only one location
			// log.Error().Msg("Error: Could not retrieve block data")
			// Creating a new error but not printing
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
	builder := string(L1Block.Header().Extra)

	// TODO(@ckartik): Move this to the chaintracer
	if st.integreationTestMode {
		builder = IntegrationTestBuilders[st.currentBlockNumberCached%int64(len(IntegrationTestBuilders))]
	}

	return blockData, builder, nil
}

type SmartContractTracerOptions struct {
	ContractClient      *rollupclient.OracleClient
	L1Client            *ethclient.Client
	StartingBlockNumber int64
	FastModeSleep       time.Duration
	NormalModeSleep     time.Duration
	FastModeSensitivity int64
	IntegrationMode     bool
}

// NewSmartContractTracer creates a new SmartContractTracer with the given options.
func NewSmartContractTracer(ctx context.Context, options SmartContractTracerOptions) Tracer {
	tracer := &SmartContractTracer{
		contractClient:       options.ContractClient,
		startingBlockNumber:  options.StartingBlockNumber,
		L1Client:             options.L1Client,
		fastModeSleep:        options.FastModeSleep,
		normalModeSleep:      options.NormalModeSleep,
		fastModeSensitivity:  options.FastModeSensitivity,
		integreationTestMode: options.IntegrationMode,
	}

	return tracer
}

type L1DataRetriver interface {
	GetTransactions(blockNumber int64) (*BlockDetails, error)
	GetWinningBuilder(blockNumber int64) (string, error)
	GetLatestBlockNumber() (int64, error)
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
