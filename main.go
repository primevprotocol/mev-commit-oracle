package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"flag"
	"fmt"
	"math/big"
	"os"
	"time"

	"database/sql"

	_ "github.com/lib/pq"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primevprotocol/mev-oracle/pkg/chaintracer"
	"github.com/primevprotocol/mev-oracle/pkg/preconf"
	"github.com/primevprotocol/mev-oracle/pkg/repository"
	"github.com/primevprotocol/mev-oracle/pkg/rollupclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// DB Setup
const (
	port = 5432          // Default port for PostgreSQL
	user = "oracle_user" // Your database user
	//  TODO(@ckartik): Move to KMS or env
	// password = "oracle_pass" // Your database password
	dbname = "oracle_db" // Your database name
)

var (
	oracleContract         = flag.String("oracle", "0x51dcB14bcb0B4eE747BE445550A4Fb53ecd31617", "Oracle contract address")
	preConfContract        = flag.String("preconf", "0xBB632720f817792578060F176694D8f7230229d9", "Preconf contract address")
	clientURL              = flag.String("rpc-url", "http://sl-bootnode:8545", "Client URL")
	l1RPCURL               = flag.String("l1-rpc-url", "http://host.docker.internal:8545", "L1 Client URL")
	privateKeyInput        = flag.String("key", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", "Private Key")
	rateLimit              = flag.Int64("rateLimit", 12, "Rate Limit in seconds")
	startBlockNumber       = flag.Int64("startBlockNumber", 0, "Start Block Number")
	onlyMonitorCommitments = flag.Bool("onlyMonitorCommitments", false, "Only monitor commitments")
	dbHost                 = flag.String("dbHost", "oracle-db", "DB Host")
	fastModeSleep          = flag.Int64("fastModeSleep", 5, "Sleep time in fast mode between data retrievials from RPC Ethereum Client")
	normalModeSleep        = flag.Int64("normalModeSleep", 12, "Sleep time in normal mode between data retrievials from RPC Ethereum Client")
	fastModeSensitivity    = flag.Int64("fastModeSensitivity", 20, "Number of blocks to be behind before fast mode is triggered")

	// TODO(@ckartik): Pull txns commited to from DB and post in Oracle payload.
	integreationTestMode = flag.Bool("integrationTestMode", false, "Integration Test Mode")

	client   *ethclient.Client
	l1Client *ethclient.Client
	pc       *preconf.Preconf
	rc       *rollupclient.OracleClient
	chainID  *big.Int
)

// Can't unittest if this isn't an interface
// Need to keep interfaces for 3rd party driveres, so you can mock for unit tests
func initDB() (db *sql.DB, err error) {

	password := os.Getenv("POSTGRES_PASSWORD")

	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		*dbHost, port, user, password, dbname)

	// Open a connection
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, err
}

type Authenticator struct {
	PrivateKey *ecdsa.PrivateKey
	ChainID    *big.Int
	Client     *ethclient.Client
}

func (a Authenticator) GetAuth() (opts *bind.TransactOpts, err error) {
	// Set transaction opts
	auth, err := bind.NewKeyedTransactorWithChainID(a.PrivateKey, a.ChainID)
	if err != nil {
		return nil, err
	}
	// Set nonce (optional)
	nonce, err := client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))

	// Set gas price (optional)
	gasPrice, err := a.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	auth.GasPrice = gasPrice

	// Set gas limit (you need to estimate or set a fixed value)
	auth.GasLimit = uint64(30000000)

	return auth, nil
}

func init() {
	var err error
	// Initialize zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(os.Stdout).With().Caller().Logger()

	/* Start of Setup */
	log.Info().Msg("Parsing flags...")
	flag.Parse()
	log.Debug().
		Str("Contract Address", *oracleContract).
		Str("Preconf Contract Address", *preConfContract).
		Str("Client URL", *clientURL).
		Str("L1 Client URL", *l1RPCURL).
		Str("Private Key", "**********").
		Int64("Rate Limit", *rateLimit).
		Int64("Start Block Number", *startBlockNumber).
		Bool("Only Monitor Commitments", *onlyMonitorCommitments).
		Msg("Flags Parsed")

	// Harder to tests with ethclient
	// it has ethclient.rpcclient
	client, err = ethclient.Dial(*clientURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the Settlement Layer client")
	}

	l1Client, err = ethclient.Dial(*l1RPCURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the L1 Ethereum client")
	}

	rc, err = rollupclient.NewOracleClient(common.HexToAddress(*oracleContract), client)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating oracle client")
	}

	pc, err = preconf.NewPreconf(common.HexToAddress(*preConfContract), client)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating preconf client")
	}

	chainID, err = client.ChainID(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting chain ID")
	}
	log.Debug().Str("Chain ID", chainID.String()).Msg("Chain ID Detected")

}

func SetBuilderMapping(authenticator Authenticator, builderName string, builderAddress string) (txnHash string, err error) {
	auth, err := authenticator.GetAuth()
	if err != nil {
		return "", err
	}

	txn, err := rc.AddBuilderAddress(auth, builderName, common.HexToAddress(builderAddress))
	if err != nil {
		return "", err
	}

	return txn.Hash().String(), nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("Error running")
	}
}

func run() (err error) {
	ctx := context.Background()
	// TODO(@ckartik): Move privatekey to AWS KMS
	privateKey, err := crypto.HexToECDSA(*privateKeyInput)
	if err != nil {
		log.Error().Err(err).Msg("Error creating private key")
		return
	}
	authenticator := Authenticator{
		PrivateKey: privateKey,
		ChainID:    chainID,
		Client:     client,
	}
	if *integreationTestMode {
		log.Info().Msg("Integration Test Mode Enabled. Setting fake builder mapping")
		for _, builder := range chaintracer.IntegrationTestBuilders {
			_, err = SetBuilderMapping(authenticator, builder, builder)
			if err != nil {
				log.Error().Err(err).Msg("Error setting builder mapping")
				return
			}
		}
	}

	db, err := initDB()
	if err != nil {
		log.Error().Err(err).Msg("Error initializing DB")
		return
	}
	txnStore := repository.NewDBTxnStore(db, pc)
	time.Sleep(5 * time.Second)
	log.Info().Msg("Sleeping to allow DB to initialize tables")

	tracer := chaintracer.NewSmartContractTracer(ctx, chaintracer.SmartContractTracerOptions{
		ContractClient:      rc,
		L1Client:            l1Client,
		StartingBlockNumber: *startBlockNumber,
		FastModeSleep:       time.Duration(*fastModeSleep),
		NormalModeSleep:     time.Duration(*normalModeSleep),
		FastModeSensitivity: *fastModeSensitivity,
		IntegrationMode:     *integreationTestMode,
	})

	for blockNumber := tracer.GetNextBlockNumber(ctx); ; blockNumber = tracer.GetNextBlockNumber(ctx) {
		log.Info().Msg("Processing")
		err = submitBlock(ctx, blockNumber, tracer, authenticator, txnStore)
		switch err { // Should be a different approach to make things cleaner
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

var (
	ErrorBlockDetails = errors.New("Error retrieving block details")
	ErrorAuth         = errors.New("Error constructing auth")

	ErrorBlockSubmission = errors.New("Error submitting block")
	ErrorUnableToFilter  = errors.New("Unable to filter transactions based on commitment")
)

type SettlementWork struct {
	commitment  repository.Commitment
	isSlash     bool
	builderName string
}

// submitblock initilaizes the retreivial and storage of commitments for a block number stored on the settlmenet layer,
// processes it with L1 block data and submits a filtered list to the settlement layer
func submitBlock(ctx context.Context, blockNumber int64, tracer chaintracer.Tracer, authenticator Authenticator, txnStore repository.CommitmentsStore) error {
	doneChan, errorChan := txnStore.UpdateCommitmentsForBlockNumber(blockNumber)
	details, builder, err := tracer.RetrieveDetails()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrorBlockDetails, err)
	}

	workChannel := settler(ctx, authenticator)

	// Get txns for the block into inclusion map
	blockTxns := make(map[string]bool)
	for _, txn := range details.Transactions {
		blockTxns[txn] = true
	}

	select {
	case <-doneChan:
		commitments, err := txnStore.RetrieveCommitments(blockNumber)
		if err != nil {
			return ErrorUnableToFilter
		}
		for _, commitment := range commitments {
			isSlash := blockTxns[commitment.TxnHash]

			workChannel <- SettlementWork{
				commitment:  commitment,
				isSlash:     isSlash,
				builderName: builder,
			}
		}
		log.Info().Msg("Received data from Store")
	case err := <-errorChan:
		return fmt.Errorf("%w: %v", ErrorUnableToFilter, err)
	}

	return nil
}

// Does the job of submitting the commitments to the rollup
// TODO(@ckartik): Optimize using Aloks method for nonce management
func settler(ctx context.Context, authenticator Authenticator) chan SettlementWork {
	workChannel := make(chan SettlementWork, 100)

	go func(ctx context.Context, workChannel <-chan SettlementWork, authenticator Authenticator) {
		for work := range workChannel {
			auth, err := authenticator.GetAuth()
			if err != nil {
				log.Fatal().Err(err).Msg("Error constructing auth")
			}
			log.Info().Int("block_number", int(work.commitment.BlockNum)).Str("txn_being_commited", work.commitment.TxnHash).Msg("Posting commitment")
			commitmentPostingTxn, err := rc.ProcessBuilderCommitmentForBlockNumber(auth, work.commitment.CommitmentIndex, big.NewInt(work.commitment.BlockNum), work.builderName, work.isSlash)
			deadlineCtx, _ := context.WithTimeout(ctx, 30*time.Second)
			reciept, err := bind.WaitMined(deadlineCtx, client, commitmentPostingTxn)
			if err != nil || reciept.Status != 1 {
				log.Error().Err(err).Msgf("Error posting commitment, receipt %v", reciept)
			}
			_ = work
			_ = auth
		}
	}(ctx, workChannel, authenticator)

	return workChannel
}
