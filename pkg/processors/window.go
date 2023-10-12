package processors

import "github.com/primevprotocol/oracle/pkg/chaintracer"

type BidSet []string

type Processor interface {
	ProcessBlock(block chaintracer.BlockDetails, openBids BidSet) (open BidSet, closed BidSet, err error)
}

type WindowAlgo struct{}

// ProcessBlock is not aware of block builder dynamics, it is simply passed a block and a set of open bids, and told to figure out which bids are closed.
func (w WindowAlgo) ProcessBlock(block chaintracer.BlockDetails, openBids BidSet) (open BidSet, closed BidSet, err error) {
	for txn := range block.Transactions {
		// Process
		_ = txn
	}
	return open, closed, nil
}
