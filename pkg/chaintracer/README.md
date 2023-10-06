## chaintracer

Chaintracer is meant to be a wrapper for any source of truth around blocks and the winning Builder. It provides a simplified interface, that retrives the winner of a block auction at the blockbuilder level and provides a method for retriving the list of txn hashes assosciated with a block.

### Trust Assumptions
We're currently piggy-backing data feeds from the following sources:
- Builder That Won: Payload.de
- Transaction List and Ordering: Infura
