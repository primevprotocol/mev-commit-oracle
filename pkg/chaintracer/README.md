## chaintracer

Chaintracer is meant to be a wrapper for any source of truth around blocks and the winning Builder. It provides a simplified interface, that retrives the winner of a block auction at the blockbuilder level and provides a method for retriving the list of txn hashes assosciated with a block.

### Interface
```go
type Tracer interface {
  // IncrementBlock moves the pointer tracing blocks forward to the next block, and returns the new current blocknumber being tracked.
	IncrementBlock() (NewBlockNumber int64)
  // RetrieveDetails retrieves the transactionHashes list and winning blocknumber for the internal blocknumber being tracked.
	RetrieveDetails() (block *BlockDetails, BlockBuilder string, err error)
}
```

### Trust Assumptions
We're currently piggy-backing data feeds from the following sources:
- Builder That Won: Payload.de
- Transaction List and Ordering: Infura
