package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"flag"
	"math/big"
	"os"
	"os/signal"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primevprotocol/oracle/pkg/chaintracer"
	"github.com/primevprotocol/oracle/pkg/rollupclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
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

var (
	contractAddress  = flag.String("contract", "0x0F81Ae3c80CD1fBa5579690Dd0425f74035DCF32", "Contract address")
	clientURL        = flag.String("rpc-url", "http://localhost:8545", "Client URL")
	privateKeyInput  = flag.String("key", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", "Private Key")
	rateLimit        = flag.Int64("rateLimit", 12, "Rate Limit in seconds")
	startBlockNumber = flag.Int64("startBlockNumber", 0, "Start Block Number")

	client  *ethclient.Client
	rc      *rollupclient.Rollupclient
	chainID *big.Int

	otelTracer             = otel.Tracer("oracle_block_submitted")
	blockSubmissionCounter metric.Int64Counter
)

func init() {
	var err error
	// Initialize zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(os.Stdout)

	/* Start of Setup */
	log.Info().Msg("Parsing flags...")
	flag.Parse()
	log.Debug().
		Str("Contract Address", *contractAddress).
		Str("Client URL", *clientURL).
		Str("Private Key", "**********").
		Int64("Rate Limit", *rateLimit).
		Int64("Start Block Number", *startBlockNumber).
		Msg("Flags Parsed")

	client, err = ethclient.Dial(*clientURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to the Ethereum client")
		panic(err)
	}

	rc, err = rollupclient.NewRollupclient(common.HexToAddress(*contractAddress), client)
	if err != nil {
		log.Error().Err(err).Msg("Error creating rollup client")
		panic(err)
	}

	chainID, err = client.ChainID(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Error getting chain ID")
		panic(err)
	}
	log.Debug().Str("Chain ID", chainID.String()).Msg("Chain ID Detected")
}

func SetBuilderMapping(pk *ecdsa.PrivateKey, builderName string, builderAddress common.Address) (txnHash string, err error) {
	auth, err := getAuth(pk, chainID, client)
	if err != nil {
		log.Error().Err(err).Msg("Failed to construct auth")
		return
	}

	txn, err := rc.AddBuilderAddress(auth, "k builder", common.HexToAddress("0x15766e4fC283Bb52C5c470648AeA2b5Ad133410a"))
	if err != nil {
		log.Error().Err(err).Msg("Error adding builder address")
		return "", err
	}

	return txn.Hash().String(), nil
}

// Have some metrics for the number of events registered
func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("Error running")
	}
}

func run() (err error) {
	// Set up OpenTelemetry.
	serviceName := "oracle"
	serviceVersion := "0.1.0"
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	otelShutdown, err := setupOTelSDK(ctx, serviceName, serviceVersion)
	if err != nil {
		return
	}

	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// TODO(@ckartik): Move privatekey to AWS KMS
	privateKey, err := crypto.HexToECDSA(*privateKeyInput)
	if err != nil {
		log.Error().Err(err).Msg("Error creating private key")
		return err
	}

	// tracer := chaintracer.NewIncrementingTracer(*startBlockNumber, time.Second*time.Duration(*rateLimit))
	tracer := chaintracer.NewSmartContractTracer(rc, ctx)
	for blockNumber := tracer.GetNextBlockNumber(ctx); ; blockNumber = tracer.GetNextBlockNumber(ctx) {
		select {
		case <-ctx.Done():
			log.Info().Msg("Shutting down")
			otelShutdown(ctx)
			log.Info().Msg("Shutdown complete")
			return nil
		default:
			log.Info().Msg("Processing")
			err = submitBlock(ctx, blockNumber, tracer, privateKey)
			switch err {
			case nil:
			case ErrorBlockDetails:
				log.Error().Err(err).Msg("Error retrieving block details")
				continue
			case ErrorAuth:
				log.Error().Err(err).Msg("Error constructing auth")
				continue
			case ErrorBlockSubmission:
				log.Error().Err(err).Msg("Error submitting block")
				continue
			default:
				log.Error().Err(err).Msg("Unknown error")
				return err
			}

		}

	}
}

var (
	ErrorBlockDetails = errors.New("Error retrieving block details")
	ErrorAuth         = errors.New("Error constructing auth")

	ErrorBlockSubmission = errors.New("Error submitting block")
)

func submitBlock(ctx context.Context, blockNumber int64, tracer chaintracer.Tracer, privateKey *ecdsa.PrivateKey) (err error) {
	ctx, span := otelTracer.Start(ctx, "block_submission")
	defer span.End()

	blockNumberAttr := attribute.Int64("blocknumber.value", blockNumber)
	span.SetAttributes(blockNumberAttr)

	details, builder, err := tracer.RetrieveDetails()
	if err != nil {
		log.Error().Int64("block_number", blockNumber).Err(err).Msg("Error retrieving block details")
		return ErrorBlockDetails
	}
	auth, err := getAuth(privateKey, chainID, client)
	if err != nil {
		log.Error().Err(err).Msg("Failed to construct auth")
		return ErrorAuth
	}
	log.Debug().Str("block_number", details.BlockNumber).Msg("Posting data to settlement layer")
	blockDataTxn, err := rc.ReceiveBlockData(auth, details.Transactions, big.NewInt(blockNumber), builder)
	if err != nil {
		log.Error().Err(err).Msg("Error posting data to settlement layer")
		return ErrorBlockSubmission
	}
	blockSubmissionCounter.Add(ctx, 1)
	log.Info().Int("data_sent_bytes", len(details.Transactions)*32).Str("txn_hash", blockDataTxn.Hash().String()).Msg("Block Data Send to Mev-Commit Settlement Contract")
	return nil
}
