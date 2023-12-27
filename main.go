package main

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
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
	rc       *rollupclient.Rollupclient
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

func getAuth(privateKey *ecdsa.PrivateKey, chainID *big.Int, client *ethclient.Client) (opts *bind.TransactOpts, err error) {
	// Set transaction opts
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
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
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	auth.GasPrice = gasPrice

	// Set gas limit (you need to estimate or set a fixed value)
	auth.GasLimit = uint64(30000000) // Example value

	return auth, nil
}

type transactor struct {
	*ethclient.Client
}

func (t *transactor) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return nil
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

	rc, err = rollupclient.NewRollupclient(common.HexToAddress(*oracleContract), &transactor{client})
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

func SetBuilderMapping(pk *ecdsa.PrivateKey, builderName string, builderAddress string) (txnHash string, err error) {
	auth, err := getAuth(pk, chainID, client)
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

	if *integreationTestMode {
		log.Info().Msg("Integration Test Mode Enabled. Setting fake builder mapping")
		for _, builder := range integrationTestBuilders {
			_, err = SetBuilderMapping(privateKey, builder, builder)
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

	if *onlyMonitorCommitments {
		for blockNumber, _ := txnStore.LargestStoredBlockNumber(); ; blockNumber++ {
			done, err := txnStore.UpdateCommitmentsForBlockNumber(int64(blockNumber))
			select {
			case <-ctx.Done():
				log.Info().Msg("Shutting down")

			case <-done:
				log.Debug().Int64("blockNumber", blockNumber).Msg("Done updating commitments")
				txns, err := txnStore.RetrieveCommitments(blockNumber)
				if err != nil {
					log.Error().Err(err).Msg("Error retrieving commitments")
					return err
				}
				for txn := range txns {
					log.Info().Str("txn", txn).Msg("Txn was commited to")
				}
			case err := <-err:
				log.Error().Err(err).Msg("Error updating commitments, skipping")
			}
		}

	}
	tracer := chaintracer.NewSmartContractTracer(ctx, chaintracer.SmartContractTracerOptions{
		ContractClient:      rc,
		L1Client:            l1Client,
		StartingBlockNumber: *startBlockNumber,
		FastModeSleep:       time.Duration(*fastModeSleep),
		NormalModeSleep:     time.Duration(*normalModeSleep),
		FastModeSensitivity: *fastModeSensitivity,
	})

	for blockNumber := tracer.GetNextBlockNumber(ctx); ; blockNumber = tracer.GetNextBlockNumber(ctx) {
		log.Info().Msg("Processing")
		err = submitBlock(ctx, blockNumber, tracer, privateKey, txnStore)
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

var integrationTestBuilders = []string{
	"0x48ddC642514370bdaFAd81C91e23759B0302C915",
	"0x972eb4Fc3c457da4C957306bE7Fa1976BB8F39A6",
	"0xA1e8FDB3bb6A0DB7aA5Db49a3512B01671686DCB",
	"0xB9286CB4782E43A202BfD426AbB72c8cb34f886c",
	"0xdaa1EEe546fc3f2d10C348d7fEfACE727C1dfa5B",
	"0x93DC0b6A7F454Dd10373f1BdA7Fe80BB549EE2F9",
	"0x426184Df456375BFfE7f53FdaF5cB48DeB3bbBE9",
	"0x41cC09BD5a97F22045fe433f1AF0B07d0AB28F58",
}

// submitblock initilaizes the retreivial and storage of commitments for a block number stored on the settlmenet layer,
// processes it with L1 block data and submits a filtered list to the settlement layer
func submitBlock(ctx context.Context, blockNumber int64, tracer chaintracer.Tracer, privateKey *ecdsa.PrivateKey, txnStore repository.CommitmentsStore) error {
	doneChan, errorChan := txnStore.UpdateCommitmentsForBlockNumber(blockNumber)
	details, builder, err := tracer.RetrieveDetails()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrorBlockDetails, err)
	}

	var transactionsToPost []string
	select {
	case <-doneChan:
		txnsCommitedTo, err := txnStore.RetrieveCommitments(blockNumber)
		if err != nil {
			return ErrorUnableToFilter
		}
		for _, txn := range details.Transactions {
			if txnsCommitedTo[txn] {
				transactionsToPost = append(transactionsToPost, txn)
			}
		}
		log.Info().Msg("Received data from Store")
	case err := <-errorChan:
		return fmt.Errorf("%w: %v", ErrorUnableToFilter, err)
	}

	log.Debug().Str("block_number", details.BlockNumber).Msg("Posting data to settlement layer")
	auth, err := getAuth(privateKey, chainID, client)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrorAuth, err)
	}

	// Have an integreation compose
	// When you have a true false branching, there is something wrong
	// It's not a useful flag
	if *integreationTestMode {
		builder = integrationTestBuilders[blockNumber%int64(len(integrationTestBuilders))]
	}

	oracleDataPostedTxn, err := rc.ReceiveBlockData(auth, transactionsToPost, big.NewInt(blockNumber), builder)
	rawTx, err := oracleDataPostedTxn.MarshalBinary()
	log.Info().Msgf("rawTxInHex: %s", common.Bytes2Hex(rawTx))
	if err != nil {
		return err
	}
	deadlineCtx, _ := context.WithTimeout(ctx, 30*time.Second)
	r, err := bind.WaitMined(deadlineCtx, client, oracleDataPostedTxn)
	log.Info().Msgf("transaction hash: %s status: %d", oracleDataPostedTxn.Hash().Hex(), r.Status)
	if err != nil {
		return err
	}
	log.Info().Int("commitment_transactions_posted", len(transactionsToPost)).Int("txns_filtered_out", len(details.Transactions)-len(transactionsToPost)).Str("submission_txn_hash", oracleDataPostedTxn.Hash().String()).Msg("Block Data Send to Mev-Commit Settlement Contract")
	return nil
}
