package l1Listener

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	blocktracker "github.com/primevprotocol/contracts-abi/clients/BlockTracker"
	"github.com/primevprotocol/mev-oracle/pkg/events"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
)

var checkInterval = 2 * time.Second

type L1Recorder interface {
	RecordL1Block(blockNum *big.Int, winner common.Address) (*types.Transaction, error)
}

type WinnerRegister interface {
	RegisterWinner(ctx context.Context, blockNum int64, winner []byte, window int64) error
}

type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

type Oracle interface {
	GetBuilder(builder string) (common.Address, error)
}

type L1Listener struct {
	logger               *slog.Logger
	l1Client             EthClient
	winnerRegister       WinnerRegister
	rollupClient         Oracle
	builderIdentityCache map[string]common.Address
	eventMgr             events.EventManager
	recorder             L1Recorder
	metrics              *metrics
}

func NewL1Listener(
	logger *slog.Logger,
	l1Client EthClient,
	winnerRegister WinnerRegister,
	oracle Oracle,
	evtMgr events.EventManager,
	recorder L1Recorder,
) *L1Listener {
	return &L1Listener{
		logger:               logger,
		l1Client:             l1Client,
		winnerRegister:       winnerRegister,
		rollupClient:         oracle,
		eventMgr:             evtMgr,
		recorder:             recorder,
		builderIdentityCache: make(map[string]common.Address),
		metrics:              newMetrics(),
	}
}

func (l *L1Listener) Metrics() []prometheus.Collector {
	return l.metrics.Collectors()
}

func (l *L1Listener) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return l.watchL1Block(egCtx)
	})

	eg.Go(func() error {
		evt := events.NewEventHandler(
			"NewL1Block",
			func(update *blocktracker.BlocktrackerNewL1Block) error {
				l.logger.Info(
					"new L1 block event",
					"block", update.BlockNumber,
					"winner", update.Winner.String(),
					"window", update.Window,
				)
				err := l.winnerRegister.RegisterWinner(
					ctx,
					update.BlockNumber.Int64(),
					update.Winner.Bytes(),
					update.Window.Int64(),
				)
				if err != nil {
					l.logger.Error(
						"failed to register winner",
						"block", update.BlockNumber,
						"winner", update.Winner.String(),
						"error", err,
					)
					return err
				}
				return nil
			},
		)

		sub, err := l.eventMgr.Subscribe(evt)
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
			l.logger.Error("L1listener error", "error", err)
		}
	}()

	return doneChan
}

func (l *L1Listener) watchL1Block(ctx context.Context) error {
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	currentBlockNo := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			blockNum, err := l.l1Client.BlockNumber(ctx)
			if err != nil {
				l.logger.Error("failed to get block number", "error", err)
				continue
			}

			if blockNum <= uint64(currentBlockNo) {
				continue
			}

			header, err := l.l1Client.HeaderByNumber(ctx, big.NewInt(int64(blockNum)))
			if err != nil {
				l.logger.Error("failed to get header", "block", blockNum, "error", err)
				continue
			}

			winner := string(bytes.ToValidUTF8(header.Extra, []byte("�")))

			builderAddr, ok := l.builderIdentityCache[winner]
			if !ok {
				builderAddr, err = l.rollupClient.GetBuilder(winner)
				if err != nil || builderAddr.Cmp(common.Address{}) == 0 {
					if errors.Is(err, ethereum.NotFound) {
						l.logger.Warn(
							"builder not registered",
							"builder", winner,
							"block", header.Number.Int64(),
						)
					}
					l.logger.Error("failed to get builder address", "error", err)
					continue
				}
				l.builderIdentityCache[winner] = builderAddr
			}

			l.logger.Info(
				"new L1 winner",
				"winner", winner,
				"block", header.Number.Int64(),
				"builder", builderAddr.String(),
			)

			winnerPostingTxn, err := l.recorder.RecordL1Block(
				big.NewInt(0).SetUint64(blockNum),
				builderAddr,
			)
			if err != nil {
				l.logger.Error("failed to register winner for block", "block", blockNum, "error", err)
				return err
			}

			l.metrics.WinnerRoundCount.WithLabelValues(builderAddr.String()).Inc()
			l.metrics.WinnerCount.Inc()

			l.logger.Info(
				"registered winner",
				"winner", builderAddr.String(),
				"block", header.Number.Int64(),
				"txn", winnerPostingTxn.Hash().String(),
			)
			currentBlockNo = int(blockNum)
		}
	}
}
