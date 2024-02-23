package l1Listener

import (
	"bytes"
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

var checkInterval = 2 * time.Second

type WinnerRegister interface {
	RegisterWinner(ctx context.Context, blockNum int64, winner string) error
}

type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
}

type L1Listener struct {
	l1Client       EthClient
	winnerRegister WinnerRegister
	metrics        *metrics
	waitOnFinality bool
}

func NewL1Listener(
	l1Client EthClient,
	winnerRegister WinnerRegister,
	waitOnFinality bool,
) *L1Listener {
	return &L1Listener{
		l1Client:       l1Client,
		winnerRegister: winnerRegister,
		metrics:        newMetrics(),
		waitOnFinality: waitOnFinality,
	}
}

func (l *L1Listener) Metrics() []prometheus.Collector {
	return l.metrics.Collectors()
}

func (l *L1Listener) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	go func() {
		defer close(doneChan)

		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()

		currentBlockNo := 0
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				blockNum, err := l.l1Client.BlockNumber(ctx)
				if err != nil {
					log.Error().Err(err).Msg("failed to get block number")
					continue
				}
				if l.waitOnFinality && blockNum > 64 {
					// This is a weak proxy for finaility to limit exposure to reorgs
					// TODO(@ckartik): Interact with Finality gadget to give a more robust implementation.
					blockNum -= 64
				}
				if blockNum <= uint64(currentBlockNo) {
					continue
				}

				header, err := l.l1Client.HeaderByNumber(ctx, big.NewInt(int64(blockNum)))
				if err != nil {
					log.Error().Err(err).
						Uint64("block", blockNum).
						Msg("failed to get header")
					continue
				}

				winner := string(bytes.ToValidUTF8(header.Extra, []byte("�")))
				if len(winner) == 0 {
					log.Warn().
						Int64("block", header.Number.Int64()).
						Msg("no winner registered")
					continue
				} else {
					err = l.winnerRegister.RegisterWinner(ctx, int64(blockNum), winner)
					if err != nil {
						log.Error().Err(err).
							Uint64("block", blockNum).
							Msg("failed to register winner for block")
						return
					}

					l.metrics.WinnerRoundCount.WithLabelValues(winner).Inc()
					l.metrics.WinnerCount.Inc()

					log.Info().
						Str("winner", winner).
						Int64("block", header.Number.Int64()).
						Msg("registered winner")
				}
				currentBlockNo = int(blockNum)
			}
		}

	}()

	return doneChan
}
