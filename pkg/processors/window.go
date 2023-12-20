package processors

import (
	"errors"

	"github.com/primevprotocol/mev-oracle/pkg/chaintracer"
)

type BidSet map[string]struct{}

type Processor interface {
	ProcessBlock(block chaintracer.BlockDetails, openBids BidSet) (open BidSet, closed BidSet, err error)
}

type WindowAlgo struct{}

// ProcessBlock is not aware of block builder dynamics, it is simply passed a block and a set of open bids, and told to figure out which bids are closed.
func (w WindowAlgo) ProcessBlock(block chaintracer.BlockDetails, openBids BidSet) (open BidSet, closed BidSet, err error) {

	if block.Transactions == nil {
		return nil, nil, errors.New("transactions cannot be nil")
	}

	// Initialize sets
	closedBids := make(BidSet)

	for _, txn := range block.Transactions {
		if _, exists := openBids[txn]; exists {
			closedBids[txn] = struct{}{}
		}
	}

	// Remove closed bids from the open set
	for txn := range closedBids {
		delete(openBids, txn)
	}

	return openBids, closedBids, nil
}
