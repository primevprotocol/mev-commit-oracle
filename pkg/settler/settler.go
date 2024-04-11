package settler

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	blocktracker "github.com/primevprotocol/contracts-abi/clients/BlockTracker"
	"github.com/primevprotocol/mev-oracle/pkg/events"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
)

type SettlementType string

const (
	SettlementTypeReward SettlementType = "reward"
	SettlementTypeSlash  SettlementType = "slash"
	SettlementTypeReturn SettlementType = "return"
)

type Settlement struct {
	CommitmentIdx   []byte
	TxHash          string
	BlockNum        int64
	Builder         []byte
	Amount          uint64
	BidID           []byte
	Type            SettlementType
	DecayPercentage int64
}

type SettlerRegister interface {
	SubscribeSettlements(ctx context.Context, window int64) <-chan Settlement
	SettlementInitiated(ctx context.Context, commitmentIdx []byte, txHash common.Hash, nonce uint64) error
}

type Oracle interface {
	ProcessBuilderCommitmentForBlockNumber(
		commitmentIdx [32]byte,
		blockNum *big.Int,
		builder string,
		isSlash bool,
		residualDecay *big.Int,
	) (*types.Transaction, error)
}

type Settler struct {
	logger          *slog.Logger
	rollupClient    Oracle
	settlerRegister SettlerRegister
	evtMgr          events.EventManager
	metrics         *metrics
}

func NewSettler(
	logger *slog.Logger,
	rollupClient Oracle,
	settlerRegister SettlerRegister,
	evtMgr events.EventManager,
) *Settler {
	return &Settler{
		logger:          logger,
		rollupClient:    rollupClient,
		settlerRegister: settlerRegister,
		evtMgr:          evtMgr,
		metrics:         newMetrics(),
	}
}

func (s *Settler) Metrics() []prometheus.Collector {
	return s.metrics.Collectors()
}

func (s *Settler) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		evt := events.NewEventHandler(
			"NewWindow",
			func(update *blocktracker.BlocktrackerNewWindow) error {
				return s.settlementExecutor(egCtx, update.Window.Int64()-2)
			},
		)

		sub, err := s.evtMgr.Subscribe(evt)
		if err != nil {
			return err
		}
		defer sub.Unsubscribe()

		select {
		case <-egCtx.Done():
			return egCtx.Err()
		case err := <-sub.Err():
			return err
		}
	})

	go func() {
		defer close(doneChan)
		if err := eg.Wait(); err != nil {
			s.logger.Error("settler error", "error", err)
		}
	}()

	return doneChan
}

func (s *Settler) settlementExecutor(ctx context.Context, window int64) error {
	cctx, unsub := context.WithCancel(ctx)
	defer unsub()

	settlementChan := s.settlerRegister.SubscribeSettlements(cctx, window)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case settlement, more := <-settlementChan:
			if !more {
				return nil
			}

			err := func() error {
				if settlement.Type == SettlementTypeReturn {
					s.logger.Warn("return settlement", "commitmentIdx", fmt.Sprintf("%x", settlement.CommitmentIdx))
					return nil
				}

				var (
					commitmentIdx [32]byte
				)

				copy(commitmentIdx[:], settlement.CommitmentIdx)

				commitmentPostingTxn, err := s.rollupClient.ProcessBuilderCommitmentForBlockNumber(
					commitmentIdx,
					big.NewInt(settlement.BlockNum),
					common.Bytes2Hex(settlement.Builder),
					settlement.Type == SettlementTypeSlash,
					big.NewInt(settlement.DecayPercentage),
				)
				if err != nil {
					return fmt.Errorf("process commitment: %w nonce %d", err, commitmentPostingTxn.Nonce())
				}

				err = s.settlerRegister.SettlementInitiated(
					ctx,
					settlement.CommitmentIdx,
					commitmentPostingTxn.Hash(),
					commitmentPostingTxn.Nonce(),
				)
				if err != nil {
					return fmt.Errorf("failed to mark settlement initiated: %w", err)
				}

				s.metrics.LastUsedNonce.Set(float64(commitmentPostingTxn.Nonce()))
				s.metrics.SettlementsPostedCount.Inc()
				s.metrics.CurrentSettlementL1Block.Set(float64(settlement.BlockNum))

				s.logger.Info(
					"builder commitment processed",
					"blockNum", settlement.BlockNum,
					"txHash", commitmentPostingTxn.Hash().Hex(),
					"builder", settlement.Builder,
					"settlementType", string(settlement.Type),
					"nonce", commitmentPostingTxn.Nonce(),
				)

				return nil
			}()
			if err != nil {
				s.logger.Error("failed to process builder commitment", "error", err)
				unsub()
				return err
			}
		}
	}
}
