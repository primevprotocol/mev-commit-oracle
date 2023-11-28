package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primevprotocol/oracle/pkg/chaintracer"
	"github.com/primevprotocol/oracle/pkg/rollupclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func getAuth(privateKey *ecdsa.PrivateKey, chainID *big.Int, client *ethclient.Client) (opts *bind.TransactOpts, err error) {
	// Set transaction opts
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to construct auth")
		return nil, err
	}
	// Set nonce (optional)
	nonce, err := client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get nonce")
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))

	// Set gas price (optional)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Failed to get gas price")
		return nil, err
	}
	auth.GasPrice = gasPrice

	// Set gas limit (you need to estimate or set a fixed value)
	auth.GasLimit = uint64(30000000) // Example value

	return auth, nil
}

type dummyTracer struct {
	blockNumberCurrent int64
}

func (d *dummyTracer) IncrementBlock() int64 {
	d.blockNumberCurrent += 1
	return d.blockNumberCurrent
}

func (d *dummyTracer) RetrieveDetails() (block *chaintracer.BlockDetails, BlockBuilder string, err error) {
	block = &chaintracer.BlockDetails{
		BlockNumber:  strconv.FormatInt(d.blockNumberCurrent, 10),
		Transactions: []string{},
	}

	for i := 0; i < 200; i++ {
		randomInt, err := rand.Int(rand.Reader, big.NewInt(1000))
		if err != nil {
			panic(err)
		}
		randomBytes := crypto.Keccak256(randomInt.Bytes())
		block.Transactions = append(block.Transactions, hex.EncodeToString(randomBytes))
	}

	sleepDuration, _ := rand.Int(rand.Reader, big.NewInt(12))
	time.Sleep(time.Duration(sleepDuration.Int64()) * time.Second)
	return block, "k builder", nil
}

func main() {
	// Initialize zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(os.Stderr)

	/* Start of Setup */
	contractAddress := flag.String("contract", "0xCf7Ed3AccA5a467e9e704C703E8D87F634fB0Fc9", "Contract address")
	clientURL := flag.String("rpc-url", "http://localhost:8545", "Client URL")
	privateKeyInput := flag.String("key", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", "Private Key")
	rateLimit := flag.Int64("rateLimit", 2, "Rate Limit in seconds")
	startBlockNumber := flag.Int64("startBlockNumber", 0, "Start Block Number")

	log.Info().Msg("Parsing flags...")
	flag.Parse()
	log.Debug().
		Str("Contract Address", *contractAddress).
		Str("Client URL", *clientURL).
		Str("Private Key", "**********").
		Int64("Rate Limit", *rateLimit).
		Int64("Start Block Number", *startBlockNumber).
		Msg("Flags Parsed")

	client, err := ethclient.Dial(*clientURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to the Ethereum client")
		return
	}

	rc, err := rollupclient.NewClient(common.HexToAddress(*contractAddress), client)
	if err != nil {
		log.Error().Err(err).Msg("Error creating rollup client")
		return
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Error getting chain ID")
		return
	}
	log.Debug().Str("Chain ID", chainID.String()).Msg("Chain ID Detected")

	// TODO(@ckartik): Move privatekey to AWS KMS
	privateKey, err := crypto.HexToECDSA(*privateKeyInput)
	if err != nil {
		log.Error().Err(err).Msg("Error creating private key")
		return
	}
	/*
		auth, err := getAuth(privateKey, chainID, client)
		if err != nil {
			log.Error().Err(err).Msg("Failed to construct auth")
			return
		}

		// Ensure we've added the builder address to the oracle.
		// An example of how to add a builder is shown below.
		//
		// txn, err := rc.AddBuilderAddress(auth, "k builder", common.HexToAddress("0x15766e4fC283Bb52C5c470648AeA2b5Ad133410a"))
		// if err != nil {
		// 	log.Error().Err(err).Msg("Error adding builder address")
		// 	return
		// }

		log.Info().Str("Transaction Hash", txn.Hash().String()).Msg("Builder Address Added")
	*/
	/* End of setup */

	// tracer := dummyTracer{10}
	tracer := chaintracer.NewIncrementingTracer(*startBlockNumber, time.Second*time.Duration(*rateLimit))
	for {
		blockNumber := tracer.IncrementBlock()
		log.Info().Int64("block_number", blockNumber).Msg("Starting to process Block")
		details, builder, err := tracer.RetrieveDetails()
		if err != nil {
			log.Error().Int64("block_number", blockNumber).Err(err).Msg("Error retrieving block details")
			continue
		}
		auth, err := getAuth(privateKey, chainID, client)
		if err != nil {
			log.Error().Err(err).Msg("Failed to construct auth")
			return
		}
		log.Debug().Str("Block Number", details.BlockNumber).Msg("Posting data to settlement layer")
		blockDataTxn, err := rc.ReceiveBlockData(auth, details.Transactions, big.NewInt(blockNumber), builder)
		if err != nil {
			log.Error().Err(err).Msg("Error posting data to sttlement layer")
			continue
		}
		log.Info().Str("Transaction Hash", blockDataTxn.Hash().String()).Msg("Block Data Send to Mev-Commit Settlement Contract")
	}
}
