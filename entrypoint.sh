#!/bin/sh

# Any environment variable setup can go here
# export MY_VAR=value

# Start the mev-oracle application
# Replace this with the actual command to start your application
./mev-oracle --rpc-url http://host.docker.internal:8545 --l1-rpc-url ${L1_URL} --startBlockNumber ${STARTING_BLOCK}

# Optional: Additional commands for logging or error handling can go here
