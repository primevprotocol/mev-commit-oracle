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
	"github.com/primevprotocol/oracle/pkg/chaintracer"
	"github.com/primevprotocol/oracle/pkg/preconf"
	"github.com/primevprotocol/oracle/pkg/repository"
	"github.com/primevprotocol/oracle/pkg/rollupclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// DB Setup
const (
	host     = "172.30.1.3"
	port     = 5432          // Default port for PostgreSQL
	user     = "oracle_user" // Your database user
	password = "oracle_pass" // Your database password
	dbname   = "oracle_db"   // Your database name
)

var (
	oracleContract         = flag.String("oracle", "0x5FC8d32690cc91D4c39d9d3abcBD16989F875707", "Oracle contract address")
	preConfContract        = flag.String("preconf", "0x9fE46736679d2D9a65F0992F2272dE9f3c7fa6e0", "Preconf contract address")
	clientURL              = flag.String("rpc-url", "http://host.docker.internal:8545", "Client URL")
	l1RPCURL               = flag.String("l1-rpc-url", "http://host.docker.internal:8545", "L1 Client URL")
	privateKeyInput        = flag.String("key", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", "Private Key")
	rateLimit              = flag.Int64("rateLimit", 12, "Rate Limit in seconds")
	startBlockNumber       = flag.Int64("startBlockNumber", 0, "Start Block Number")
	onlyMonitorCommitments = flag.Bool("onlyMonitorCommitments", false, "Only monitor commitments")

	// TODO(@ckartik): Pull txns commited to from DB and post in Oracle payload.
	integreationTestMode = flag.Bool("integrationTestMode", false, "Integration Test Mode")

	client   *ethclient.Client
	l1Client *ethclient.Client
	pc       *preconf.PreConfClient
	rc       *rollupclient.Rollupclient
	chainID  *big.Int
)

func initDB() (db *sql.DB, err error) {
	// Connection string
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open a connection
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Error().Err(err).Msg("Error opening DB")
		return nil, err
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		log.Error().Err(err).Msg("Error pinging DB")
	}

	return db, err
}

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

func init() {
	var err error
	// Initialize zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(os.Stdout)

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

	client, err = ethclient.Dial(*clientURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to the Ethereum client")
		os.Exit(1)
	}

	l1Client, err = ethclient.Dial(*l1RPCURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to the Ethereum client")
		os.Exit(1)
	}

	rc, err = rollupclient.NewRollupclient(common.HexToAddress(*oracleContract), client)
	if err != nil {
		log.Error().Err(err).Msg("Error creating oracle client")
		os.Exit(1)
	}

	pc, err = preconf.NewPreConfClient(common.HexToAddress(*preConfContract), client)
	if err != nil {
		log.Error().Err(err).Msg("Error creating preconf client")
		os.Exit(1)
	}

	chainID, err = client.ChainID(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("Error getting chain ID")
		os.Exit(1)
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

	// txnFilter := chaintracer.NewTransactionCommitmentFilter(pc)

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
	tracer := chaintracer.NewSmartContractTracer(ctx, rc, l1Client, *startBlockNumber)
	for blockNumber := tracer.GetNextBlockNumber(ctx); ; blockNumber = tracer.GetNextBlockNumber(ctx) {
		log.Info().Msg("Processing")
		err = submitBlock(ctx, blockNumber, tracer, privateKey, txnStore)
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

var (
	ErrorBlockDetails = errors.New("Error retrieving block details")
	ErrorAuth         = errors.New("Error constructing auth")

	ErrorBlockSubmission = errors.New("Error submitting block")
	ErrorUnableToFilter  = errors.New("Unable to filter transactions based on commitment")
)

// submitblock initilaizes the retreivial and storage of commitments for a block number stored on the settlmenet layer,
// processes it with L1 block data and submits a filtered list to the settlement layer
func submitBlock(ctx context.Context, blockNumber int64, tracer chaintracer.Tracer, privateKey *ecdsa.PrivateKey, txnStore repository.CommitmentsStore) (err error) {
	doneChan, errorChan := txnStore.UpdateCommitmentsForBlockNumber(blockNumber)
	details, builder, err := tracer.RetrieveDetails()
	if err != nil {
		log.Error().Int64("block_number", blockNumber).Err(err).Msg("Error retrieving block details")
		return ErrorBlockDetails
	}

	var transactionsToPost []string
	select {
	case <-doneChan:
		txnsCommitedTo, err := txnStore.RetrieveCommitments(blockNumber)
		if err != nil {
			log.Error().Err(err).Msg("Error retrieving commitments")
			return ErrorUnableToFilter
		}
		for _, txn := range details.Transactions {
			if txnsCommitedTo[txn] {
				transactionsToPost = append(transactionsToPost, txn)
			}
		}
		log.Info().Msg("Received data from Store")
	case err := <-errorChan:
		log.Error().Err(err).Msg("Error from Store")
		return ErrorUnableToFilter
	}

	log.Debug().Str("block_number", details.BlockNumber).Msg("Posting data to settlement layer")
	auth, err := getAuth(privateKey, chainID, client)
	if err != nil {
		log.Error().Err(err).Msg("Failed to construct auth")
		return ErrorAuth
	}

	oracleDataPostedTxn, err := rc.ReceiveBlockData(auth, transactionsToPost, big.NewInt(blockNumber), builder)
	if err != nil {
		log.Error().Err(err).Msg("Error posting data to settlement layer")
		return ErrorBlockSubmission
	}

	log.Info().Int("commitment_transactions_posted", len(transactionsToPost)).Int("txns_filtered_out", len(details.Transactions)-len(transactionsToPost)).Str("submission_txn_hash", oracleDataPostedTxn.Hash().String()).Msg("Block Data Send to Mev-Commit Settlement Contract")
	return nil
}
