package l1Listener

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

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
}

func NewL1Listener(
	l1Client EthClient,
	winnerRegister WinnerRegister,
) *L1Listener {
	return &L1Listener{
		l1Client:       l1Client,
		winnerRegister: winnerRegister,
	}
}

func (l *L1Listener) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	go func() {
		defer close(doneChan)

		ticker := time.NewTicker(2 * time.Second)
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

				if blockNum <= uint64(currentBlockNo) {
					continue
				}

				currentBlockNo = int(blockNum)
				header, err := l.l1Client.HeaderByNumber(ctx, big.NewInt(int64(currentBlockNo)))
				if err != nil {
					log.Error().Err(err).Msg("failed to get header")
					continue
				}

				winner := string(header.Extra)
				err = l.winnerRegister.RegisterWinner(ctx, int64(currentBlockNo), winner)
				if err != nil {
					log.Error().Err(err).
						Int64("block", int64(currentBlockNo)).
						Msg("failed to register winner for block")
					return
				}

				log.Info().
					Str("winner", winner).
					Int64("block", header.Number.Int64()).
					Msg("registered winner")
			}
		}

	}()

	return doneChan
}
