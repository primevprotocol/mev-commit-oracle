package processors

import "github.com/primevprotocol/oracle/pkg/chaintracer"

type BidSet []string

type Processor interface {
	ProcessBlock(block chaintracer.Block, openBids BidSet) (open BidSet, closed BidSet, err error)
}

type WindowAlgo struct{}

func (w WindowAlgo) ProcessBlock(block chaintracer.Block, openBids BidSet) (open BidSet, closed BidSet, err error) {

}
