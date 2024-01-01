package node

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primevprotocol/mev-oracle/pkg/l1Listener"
	"github.com/primevprotocol/mev-oracle/pkg/preconf"
	"github.com/primevprotocol/mev-oracle/pkg/rollupclient"
	"github.com/primevprotocol/mev-oracle/pkg/settler"
	"github.com/primevprotocol/mev-oracle/pkg/store"
	"github.com/primevprotocol/mev-oracle/pkg/updater"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/sha3"
)

type Options struct {
	PrivateKey          *ecdsa.PrivateKey
	HTTPPort            int
	SettlementRPCUrl    string
	L1RPCUrl            string
	OracleContractAddr  common.Address
	PreconfContractAddr common.Address
	PgHost              string
	PgPort              int
	PgUser              string
	PgPassword          string
	PgDbname            string
	LaggerdMode         int
	OverrideWinners     []string
}

type Node struct {
	waitClose func()
	dbCloser  io.Closer
}

func NewNode(opts *Options) (*Node, error) {
	nd := &Node{}

	db, err := initDB(opts)
	if err != nil {
		log.Error().Err(err).Msg("failed initializing DB")
		return nil, err
	}
	nd.dbCloser = db

	st, err := store.NewStore(db)
	if err != nil {
		log.Error().Err(err).Msg("failed initializing store")
		return nil, err
	}

	owner := getEthAddressFromPubKey(opts.PrivateKey.Public().(*ecdsa.PublicKey))

	settlementClient, err := ethclient.Dial(opts.SettlementRPCUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to the settlement layer")
		return nil, err
	}

	chainID, err := settlementClient.ChainID(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed getting chain ID")
		return nil, err
	}

	l1Client, err := ethclient.Dial(opts.L1RPCUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to the L1 Ethereum client")
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	var listenerL1Client l1Listener.EthClient

	listenerL1Client = l1Client
	if opts.LaggerdMode > 0 {
		listenerL1Client = &laggerdL1Client{EthClient: listenerL1Client, amount: opts.LaggerdMode}
	}

	preconfContract, err := preconf.NewPreconfCaller(opts.PreconfContractAddr, settlementClient)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to instantiate preconf contract")
		return nil, err
	}

	oracleContract, err := rollupclient.NewOracleClient(opts.OracleContractAddr, settlementClient)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to instantiate oracle contract")
		return nil, err
	}

	if opts.OverrideWinners != nil && len(opts.OverrideWinners) > 0 {
		listenerL1Client = &winnerOverrideL1Client{EthClient: listenerL1Client, winners: opts.OverrideWinners}
		for _, winner := range opts.OverrideWinners {
			err := setBuilderMapping(
				ctx,
				opts.PrivateKey,
				chainID,
				settlementClient,
				oracleContract,
				winner,
				winner,
			)
			if err != nil {
				log.Fatal().Err(err).Msg("failed to set builder mapping")
				return nil, err
			}
		}
	}

	l1Lis := l1Listener.NewL1Listener(listenerL1Client, st)
	l1LisClosed := l1Lis.Start(ctx)

	callOpts := bind.CallOpts{
		Pending: false,
		From:    owner,
		Context: ctx,
	}

	pc := &preconf.PreconfCallerSession{Contract: preconfContract, CallOpts: callOpts}
	oc := &rollupclient.OracleClientSession{Contract: oracleContract, CallOpts: callOpts}

	updtr := updater.NewUpdater(l1Client, st, oc, pc)
	updtrClosed := updtr.Start(ctx)

	settlr := settler.NewSettler(
		opts.PrivateKey,
		chainID,
		owner,
		oracleContract,
		st,
		settlementClient,
	)
	settlrClosed := settlr.Start(ctx)

	nd.waitClose = func() {
		cancel()

		closeChan := make(chan struct{})
		go func() {
			defer close(closeChan)

			<-l1LisClosed
			<-updtrClosed
			<-settlrClosed
		}()

		<-closeChan
	}

	return nd, nil
}

func (n *Node) Close() (err error) {
	defer func() {
		if n.dbCloser != nil {
			if err2 := n.dbCloser.Close(); err2 != nil {
				err = errors.Join(err, err2)
			}
		}
	}()
	workersClosed := make(chan struct{})
	go func() {
		defer close(workersClosed)

		if n.waitClose != nil {
			n.waitClose()
		}
	}()

	select {
	case <-workersClosed:
		log.Info().Msg("All workers closed")
		return nil
	case <-time.After(10 * time.Second):
		log.Error().Msg("Timeout waiting for workers to close")
		return errors.New("timeout waiting for workers to close")
	}
}

func initDB(opts *Options) (db *sql.DB, err error) {
	// Connection string
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		opts.PgHost, opts.PgPort, opts.PgUser, opts.PgPassword, opts.PgDbname,
	)

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

func getEthAddressFromPubKey(key *ecdsa.PublicKey) common.Address {
	pbBytes := crypto.FromECDSAPub(key)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(pbBytes[1:])
	address := hash.Sum(nil)[12:]

	return common.BytesToAddress(address)
}

type laggerdL1Client struct {
	l1Listener.EthClient
	amount int
}

func (l *laggerdL1Client) BlockNumber(ctx context.Context) (uint64, error) {
	blkNum, err := l.EthClient.BlockNumber(ctx)
	if err != nil {
		return 0, err
	}

	return blkNum - uint64(l.amount), nil
}

type winnerOverrideL1Client struct {
	l1Listener.EthClient
	winners []string
}

func (w *winnerOverrideL1Client) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	hdr, err := w.EthClient.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	idx := number.Int64() % int64(len(w.winners))
	hdr.Extra = []byte(w.winners[idx])

	return hdr, nil
}

func setBuilderMapping(
	ctx context.Context,
	privateKey *ecdsa.PrivateKey,
	chainID *big.Int,
	client *ethclient.Client,
	rc *rollupclient.OracleClient,
	builderName string,
	builderAddress string,
) error {
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return err
	}
	nonce, err := client.PendingNonceAt(ctx, auth.From)
	if err != nil {
		return err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	gasTip, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		return err
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return err
	}

	gasFeeCap := new(big.Int).Add(gasTip, gasPrice)

	auth.GasFeeCap = gasFeeCap
	auth.GasTipCap = gasTip

	txn, err := rc.AddBuilderAddress(auth, builderName, common.HexToAddress(builderAddress))
	if err != nil {
		return err
	}

	_, err = bind.WaitMined(ctx, client, txn)
	if err != nil {
		return err
	}

	return nil
}
