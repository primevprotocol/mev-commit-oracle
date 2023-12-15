package chaintracer

import (
	"context"
	"math/big"

	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
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

	// TODO(@ckartik): Use stored block number on contract instead of incrementing (For Failure reslience)
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
