package updater

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	preconf "github.com/primevprotocol/contracts-abi/clients/PreConfCommitmentStore"
	"github.com/primevprotocol/mev-oracle/pkg/settler"
	"github.com/prometheus/client_golang/prometheus"
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
		amount uint64,
		builder string,
		bidID []byte,
		settlementType settler.SettlementType,
	) error
}

type L1Client interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
}

type Oracle interface {
	GetBuilder(builder string) (common.Address, error)
}

type Preconf interface {
	GetCommitmentsByBlockNumber(blockNum *big.Int) ([][32]byte, error)
	GetCommitment(commitmentIdx [32]byte) (preconf.PreConfCommitmentStorePreConfCommitment, error)
}

type Updater struct {
	l1Client             L1Client
	winnerRegister       WinnerRegister
	preconfClient        Preconf
	rollupClient         Oracle
	builderIdentityCache map[string]common.Address
	metrics              *metrics
}

func NewUpdater(
	l1Client L1Client,
	winnerRegister WinnerRegister,
	rollupClient Oracle,
	preconfClient Preconf,
) *Updater {
	return &Updater{
		l1Client:             l1Client,
		winnerRegister:       winnerRegister,
		preconfClient:        preconfClient,
		rollupClient:         rollupClient,
		builderIdentityCache: make(map[string]common.Address),
		metrics:              newMetrics(),
	}
}

func (u *Updater) Metrics() []prometheus.Collector {
	return u.metrics.Collectors()
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
				unsub()
				return
			case winner, more := <-winnerChan:
				if !more {
					unsub()
					goto RESTART
				}
				u.metrics.UpdaterTriggerCount.Inc()

				err := func() error {
					var err error
					builderAddr, ok := u.builderIdentityCache[winner.Winner]
					if !ok {
						builderAddr, err = u.rollupClient.GetBuilder(winner.Winner)
						if err != nil {
							if errors.Is(err, ethereum.NotFound) {
								log.Warn().
									Str("builder", winner.Winner).
									Msg("builder not registered")
								return u.winnerRegister.UpdateComplete(ctx, winner.BlockNumber)
							}
							return fmt.Errorf("failed to get builder address: %w", err)
						}
						u.builderIdentityCache[winner.Winner] = builderAddr
					}

					blk, err := u.l1Client.BlockByNumber(ctx, big.NewInt(winner.BlockNumber))
					if err != nil {
						return fmt.Errorf("failed to get block by number: %w", err)
					}

					txnsInBlock := make(map[string]int)
					for posInBlock, tx := range blk.Transactions() {
						txnsInBlock[strings.TrimPrefix(tx.Hash().Hex(), "0x")] = posInBlock
					}

					commitmentIndexes, err := u.preconfClient.GetCommitmentsByBlockNumber(
						big.NewInt(winner.BlockNumber),
					)
					if err != nil {
						return fmt.Errorf("failed to get commitments by block number: %w", err)
					}

					log.Debug().
						Int("commitments_count", len(commitmentIndexes)).
						Int("txns_count", len(txnsInBlock)).
						Int64("blockNumber", winner.BlockNumber).
						Msg("commitment indexes")

					total, rewards, slashes := 0, 0, 0
					for _, index := range commitmentIndexes {
						commitment, err := u.preconfClient.GetCommitment(index)
						if err != nil {
							return fmt.Errorf("failed to get commitment: %w", err)
						}

						settlementType := settler.SettlementTypeReturn

						if commitment.Commiter.Cmp(builderAddr) == 0 {
							commitmentTxnHashes := strings.Split(commitment.TxnHash, ",")
							settlementType = settler.SettlementTypeReward

							// Ensure Bundle is atomic and present in the block
							for i := 0; i < len(commitmentTxnHashes); i++ {
								posInBlock, found := txnsInBlock[commitmentTxnHashes[i]]
								if !found || posInBlock != txnsInBlock[commitmentTxnHashes[0]]+i {
									settlementType = settler.SettlementTypeSlash
									break
								}
							}
						}

						err = u.winnerRegister.AddSettlement(
							ctx,
							index[:],
							commitment.TxnHash,
							winner.BlockNumber,
							commitment.Bid,
							commitment.Commiter.Hex(),
							commitment.CommitmentHash[:],
							settlementType,
						)
						if err != nil {
							return fmt.Errorf("failed to add settlement: %w", err)
						}

						total++
						switch settlementType {
						case settler.SettlementTypeSlash:
							slashes++
						case settler.SettlementTypeReward:
							rewards++
						}
					}

					err = u.winnerRegister.UpdateComplete(ctx, winner.BlockNumber)
					if err != nil {
						return fmt.Errorf("failed to update completion of block updates: %w", err)
					}

					u.metrics.CommitmentsCount.Add(float64(total))
					u.metrics.RewardsCount.Add(float64(rewards))
					u.metrics.SlashesCount.Add(float64(slashes))
					u.metrics.BlockCommitmentsCount.Inc()

					log.Info().
						Int("total", total).
						Int("rewards", rewards).
						Int("slashes", slashes).
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
