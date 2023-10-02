package chaintracer

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type TransactionHash []byte

// TODO(@ckartik): Replace with cappella.Block when I have an interenet connection.
// We want to define the structure of a block in a unique location.
type Block struct {
	txns []TransactionHash
}

type Winner struct {
	blockNumber big.Int
	slotNumber  big.Int
	epoch       big.Int

	pubkey  ecdsa.PubKey
	address common.Address
}

type Tracer interface {
	GetLatestWinner() (winningBuilder Winner, err error)
	GetWinnerAt(blocknumber *big.Int) (winningBuilder Winner, err error)
	GetBlockAt(blocknumber *big.Int) (block capella.Block, err error)
}
