package updater

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/primevprotocol/mev-oracle/pkg/preconf"
	"github.com/primevprotocol/mev-oracle/pkg/rollupclient"
	"github.com/rs/zerolog/log"
)

type BlockWinner struct {
	BlockNumber int64
	Winner      string
}

type WinnerRegister interface {
	SubscribeWinners(ctx context.Context) <-chan BlockWinner
	UpdateComplete(ctx context.Context, blockNum int64) error
	AddSettlement(
		ctx context.Context,
		commitmentIdx []byte,
		txHash string,
		blockNum int64,
		builder string,
		isSlash bool,
	) error
}

type Updater struct {
	owner                common.Address
	l1Client             *ethclient.Client
	winnerRegister       WinnerRegister
	preconfClient        *preconf.Preconf
	rollupClient         *rollupclient.OracleClient
	builderIdentityCache map[string]common.Address
}

func NewUpdater(
	owner common.Address,
	l1Client *ethclient.Client,
	winnerRegister WinnerRegister,
	preconfClient *preconf.Preconf,
	rollupClient *rollupclient.OracleClient,
) *Updater {
	return &Updater{
		owner:                owner,
		l1Client:             l1Client,
		winnerRegister:       winnerRegister,
		preconfClient:        preconfClient,
		rollupClient:         rollupClient,
		builderIdentityCache: make(map[string]common.Address),
	}
}

func (u *Updater) getCallOpts(ctx context.Context) *bind.CallOpts {
	return &bind.CallOpts{
		Pending: false,
		From:    u.owner,
		Context: ctx,
	}
}

func (u *Updater) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	go func() {
		defer close(doneChan)

	RESTART:
		cctx, unsub := context.WithCancel(ctx)
		winnerChan := u.winnerRegister.SubscribeWinners(cctx)

		for {
			select {
			case <-ctx.Done():
				return
			case winner, more := <-winnerChan:
				if !more {
					unsub()
					goto RESTART
				}

				err := func() error {
					var err error
					builderAddr, ok := u.builderIdentityCache[winner.Winner]
					if !ok {
						builderAddr, err = u.rollupClient.GetBuilder(u.getCallOpts(ctx), winner.Winner)
						if err != nil {
							return fmt.Errorf("failed to get builder address: %w", err)
						}
						u.builderIdentityCache[winner.Winner] = builderAddr
					}

					blk, err := u.l1Client.BlockByNumber(ctx, big.NewInt(winner.BlockNumber))
					if err != nil {
						return fmt.Errorf("failed to get block by number: %w", err)
					}

					txnsInBlock := make(map[string]struct{})
					for _, tx := range blk.Transactions() {
						txnsInBlock[tx.Hash().Hex()] = struct{}{}
					}

					commitmentIndexes, err := u.preconfClient.GetCommitmentsByBlockNumber(
						u.getCallOpts(ctx),
						big.NewInt(winner.BlockNumber),
					)
					if err != nil {
						return fmt.Errorf("failed to get commitments by block number: %w", err)
					}

					count := 0
					for _, index := range commitmentIndexes {
						commitment, err := u.preconfClient.GetCommitment(u.getCallOpts(ctx), index)
						if err != nil {
							return fmt.Errorf("failed to get commitment: %w", err)
						}

						if commitment.Commiter.Cmp(builderAddr) == 0 {
							_, ok := txnsInBlock[commitment.TxnHash]
							err = u.winnerRegister.AddSettlement(
								ctx,
								index[:],
								commitment.TxnHash,
								winner.BlockNumber,
								winner.Winner,
								!ok,
							)
							if err != nil {
								return fmt.Errorf("failed to add settlement: %w", err)
							}
							count++
						}
					}

					err = u.winnerRegister.UpdateComplete(ctx, winner.BlockNumber)
					if err != nil {
						return fmt.Errorf("failed to update completion of block updates: %w", err)
					}

					log.Info().
						Int("count", count).
						Int64("blockNumber", winner.BlockNumber).
						Str("winner", winner.Winner).
						Msg("added settlements")

					return nil
				}()

				if err != nil {
					log.Error().Err(err).
						Int64("blockNumber", winner.BlockNumber).
						Str("winner", winner.Winner).
						Msg("failed to process settlements")
					unsub()
					goto RESTART
				}
			}
		}
	}()

	return doneChan
}
