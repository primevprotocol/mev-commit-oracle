package settler

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primevprotocol/mev-oracle/pkg/keysigner"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

var (
	allowedPendingTxnCount = 128
)

type Settlement struct {
	CommitmentIdx []byte
	TxHash        string
	BlockNum      int64
	Builder       string
	IsSlash       bool
}

type SettlerRegister interface {
	LastNonce() (int64, error)
	PendingTxnCount() (int, error)
	SubscribeSettlements(ctx context.Context) <-chan Settlement
	SettlementInitiated(ctx context.Context, commitmentIdx []byte, txHash common.Hash, nonce uint64) error
	MarkSettlementComplete(ctx context.Context, nonce uint64) (int, error)
}

type Oracle interface {
	ProcessBuilderCommitmentForBlockNumber(
		opts *bind.TransactOpts,
		commitmentIdx [32]byte,
		blockNum *big.Int,
		builder string,
		isSlash bool,
	) (*types.Transaction, error)
}

type Transactor interface {
	bind.ContractTransactor
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	BlockNumber(ctx context.Context) (uint64, error)
}

type Settler struct {
	keySigner       keysigner.KeySigner
	chainID         *big.Int
	owner           common.Address
	rollupClient    Oracle
	settlerRegister SettlerRegister
	client          Transactor
	metrics         *metrics
}

func NewSettler(
	keySigner keysigner.KeySigner,
	chainID *big.Int,
	owner common.Address,
	rollupClient Oracle,
	settlerRegister SettlerRegister,
	client Transactor,
) *Settler {
	return &Settler{
		rollupClient:    rollupClient,
		settlerRegister: settlerRegister,
		owner:           owner,
		client:          client,
		keySigner:       keySigner,
		chainID:         chainID,
		metrics:         newMetrics(),
	}
}

func (s *Settler) getTransactOpts(ctx context.Context) (*bind.TransactOpts, error) {
	auth, err := s.keySigner.GetAuth(s.chainID)
	if err != nil {
		return nil, err
	}
	nonce, err := s.client.PendingNonceAt(ctx, auth.From)
	if err != nil {
		return nil, err
	}
	usedNonce, err := s.settlerRegister.LastNonce()
	if err != nil {
		return nil, err
	}
	if nonce <= uint64(usedNonce) {
		nonce = uint64(usedNonce + 1)
	}
	auth.Nonce = big.NewInt(int64(nonce))

	gasTip, err := s.client.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, err
	}

	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	gasFeeCap := new(big.Int).Add(gasTip, gasPrice)

	auth.GasFeeCap = gasFeeCap
	auth.GasTipCap = gasTip

	return auth, nil
}

func (s *Settler) Metrics() []prometheus.Collector {
	return s.metrics.Collectors()
}

func (s *Settler) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	go func() {
		queryTicker := time.NewTicker(500 * time.Millisecond)
		defer queryTicker.Stop()

		lastBlock := uint64(0)
		for {
			select {
			case <-ctx.Done():
				return
			case <-queryTicker.C:
			}

			currentBlock, err := s.client.BlockNumber(ctx)
			if err != nil {
				log.Error().Err(err).Msg("failed to get block number")
				continue
			}

			if currentBlock <= lastBlock {
				continue
			}

			lastNonce, err := s.client.NonceAt(
				ctx,
				s.owner,
				new(big.Int).SetUint64(currentBlock),
			)
			if err != nil {
				log.Error().Err(err).Msg("failed to get nonce")
				continue
			}

			count, err := s.settlerRegister.MarkSettlementComplete(ctx, lastNonce)
			if err != nil {
				log.Error().Err(err).Msg("failed to mark settlement complete")
				continue
			}

			s.metrics.LastConfirmedNonce.Set(float64(lastNonce))
			s.metrics.LastConfirmedBlock.Set(float64(currentBlock))
			s.metrics.SettlementsConfirmedCount.Add(float64(count))

			if count > 0 {
				log.Info().Int("count", count).Msg("marked settlement complete")
			}

			lastBlock = currentBlock
		}

	}()

	go func() {
		defer close(doneChan)

	RESTART:
		cctx, unsub := context.WithCancel(ctx)
		settlementChan := s.settlerRegister.SubscribeSettlements(cctx)

		for {
			select {
			case <-ctx.Done():
				unsub()
				return
			case settlement, more := <-settlementChan:
				if !more {
					unsub()
					goto RESTART
				}

				err := func() error {
					pendingTxns, err := s.settlerRegister.PendingTxnCount()
					if err != nil {
						return err
					}

					if pendingTxns > allowedPendingTxnCount {
						time.Sleep(5 * time.Second)
						return errors.New("too many pending txns")
					}

					opts, err := s.getTransactOpts(ctx)
					if err != nil {
						return err
					}

					var commitmentIdx [32]byte
					copy(commitmentIdx[:], settlement.CommitmentIdx)

					commitmentPostingTxn, err := s.rollupClient.ProcessBuilderCommitmentForBlockNumber(
						opts,
						commitmentIdx,
						big.NewInt(settlement.BlockNum),
						settlement.Builder,
						settlement.IsSlash,
					)
					if err != nil {
						return err
					}

					err = s.settlerRegister.SettlementInitiated(
						ctx,
						settlement.CommitmentIdx,
						commitmentPostingTxn.Hash(),
						commitmentPostingTxn.Nonce(),
					)
					if err != nil {
						return err
					}

					s.metrics.LastUsedNonce.Set(float64(commitmentPostingTxn.Nonce()))
					s.metrics.SettlementsPostedCount.Inc()
					s.metrics.CurrentSettlementL1Block.Set(float64(settlement.BlockNum))

					log.Info().
						Int64("blockNum", settlement.BlockNum).
						Str("txHash", commitmentPostingTxn.Hash().Hex()).
						Str("builder", settlement.Builder).
						Bool("isSlash", settlement.IsSlash).
						Msg("builder commitment processed")

					return nil
				}()
				if err != nil {
					log.Error().Err(err).Msg("failed to process builder commitment")
					unsub()
					goto RESTART
				}
			}
		}
	}()

	return doneChan
}
