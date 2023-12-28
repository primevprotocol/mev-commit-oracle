package l1Listener

import (
	"context"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

type WinnerRegister interface {
	RegisterWinner(ctx context.Context, blockNum int64, winner string) error
}

type EthClient interface {
	SubscribeNewHead(ctx context.Context, ch chan<- *types.Header) (ethereum.Subscription, error)
}

type L1Listener struct {
	l1Client       EthClient
	winnerRegister WinnerRegister
	selector       func(header *types.Header) string
}

func NewL1Listener(
	l1Client EthClient,
	winnerRegister WinnerRegister,
	selector func(header *types.Header) string,
) *L1Listener {
	return &L1Listener{
		l1Client:       l1Client,
		winnerRegister: winnerRegister,
		selector:       selector,
	}
}

func (l *L1Listener) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	go func() {
		defer close(doneChan)

		for {
			headerChan := make(chan *types.Header)
			sub, err := l.l1Client.SubscribeNewHead(ctx, headerChan)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					return
				}
				time.Sleep(5 * time.Second)
				continue
			}
			defer sub.Unsubscribe()

			for {
				select {
				case <-ctx.Done():
					return
				case header := <-headerChan:
					winner := string(header.Extra)
					if l.selector != nil {
						winner = l.selector(header)
					}
					err := l.winnerRegister.RegisterWinner(ctx, header.Number.Int64(), winner)
					if err != nil {
						log.Error().Err(err).
							Int64("block", header.Number.Int64()).
							Msg("failed to register winner for block")
						return
					}

					log.Info().
						Str("winner", winner).
						Int64("block", header.Number.Int64()).
						Msg("registered winner")
				}
			}
		}
	}()

	return doneChan
}
