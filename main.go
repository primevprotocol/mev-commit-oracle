package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"math/big"
	"os"
	"sync"

	"database/sql"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/sha3"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primevprotocol/mev-oracle/pkg/chaintracer"
	"github.com/primevprotocol/mev-oracle/pkg/l1Listener"
	"github.com/primevprotocol/mev-oracle/pkg/preconf"
	"github.com/primevprotocol/mev-oracle/pkg/rollupclient"
	"github.com/primevprotocol/mev-oracle/pkg/settler"
	"github.com/primevprotocol/mev-oracle/pkg/store"
	"github.com/primevprotocol/mev-oracle/pkg/updater"
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
	oracleContract  = flag.String("oracle", "0xA8Efc1287bAEbbD19052CAF62F265E668fcF2146", "Oracle contract address")
	preConfContract = flag.String("preconf", "0xBB632720f817792578060F176694D8f7230229d9", "Preconf contract address")
	clientURL       = flag.String("rpc-url", "http://sl-bootnode:8545", "Client URL")
	l1RPCURL        = flag.String("l1-rpc-url", "http://host.docker.internal:8545", "L1 Client URL")
	privateKeyInput = flag.String("key", "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", "Private Key")
	dbHost          = flag.String("dbHost", "oracle-db", "DB Host")
	// TODO(@ckartik): Pull txns commited to from DB and post in Oracle payload.
	integrationTestMode = flag.Bool("integrationTestMode", false, "Integration Test Mode")

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
	lock       sync.Mutex
}

func (a Authenticator) GetAuth() (opts *bind.TransactOpts, err error) {
	// Set transaction opts
	auth, err := bind.NewKeyedTransactorWithChainID(a.PrivateKey, a.ChainID)
	if err != nil {
		return nil, err
	}
	// Set nonce (optional)
	// TODO(@ckartik): This should really only be needed once
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
	auth.GasPrice = gasPrice.Mul(gasPrice, big.NewInt(4))

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

func SetBuilderMapping(
	authenticator Authenticator,
	builderName string,
	builderAddress string,
) (txnHash string, err error) {
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

	if *integrationTestMode {
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

	st, err := store.NewStore(db)
	if err != nil {
		log.Error().Err(err).Msg("Error initializing store")
		return
	}

	owner := getEthAddressFromPubKey(privateKey.Public().(*ecdsa.PublicKey))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var selector func(*types.Header) string
	if *integrationTestMode {
		selector = func(header *types.Header) string {
			idx := header.Number.Int64() % int64(len(chaintracer.IntegrationTestBuilders))
			return chaintracer.IntegrationTestBuilders[idx]
		}
	}

	l1Lis := l1Listener.NewL1Listener(l1Client, st, selector)
	updtr := updater.NewUpdater(owner, l1Client, st, pc, rc)
	settlr := settler.NewSettler(rc, st, owner, client, privateKey, chainID)

	l1LisClosed := l1Lis.Start(ctx)
	updtrClosed := updtr.Start(ctx)
	settlrClosed := settlr.Start(ctx)

	select {
	case <-l1LisClosed:
		log.Error().Msg("L1 Listener closed")
	case <-updtrClosed:
		log.Error().Msg("Updater closed")
	case <-settlrClosed:
		log.Error().Msg("Settler closed")
	}

	return
}

func getEthAddressFromPubKey(key *ecdsa.PublicKey) common.Address {
	pbBytes := crypto.FromECDSAPub(key)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pbBytes[1:])
	address := hash.Sum(nil)[12:]

	return common.BytesToAddress(address)
}
