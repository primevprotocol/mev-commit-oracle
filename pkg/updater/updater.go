package updater

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	lru "github.com/hashicorp/golang-lru/v2"
	preconf "github.com/primevprotocol/contracts-abi/clients/PreConfCommitmentStore"
	"github.com/primevprotocol/mev-oracle/pkg/events"
	"github.com/primevprotocol/mev-oracle/pkg/settler"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sync/errgroup"
)

type Winner struct {
	Winner []byte
	Window int64
}

type WinnerRegister interface {
	AddEncryptedCommitment(
		ctx context.Context,
		commitmentIdx []byte,
		committer []byte,
		commitmentHash []byte,
		commitmentSignature []byte,
		blockNum int64,
	) error
	IsSettled(ctx context.Context, commitmentIdx []byte) (bool, error)
	GetWinner(ctx context.Context, blockNum int64) (Winner, error)
	AddSettlement(
		ctx context.Context,
		commitmentIdx []byte,
		txHash string,
		blockNum int64,
		amount uint64,
		builder []byte,
		bidID []byte,
		settlementType settler.SettlementType,
		decayPercentage int64,
		window int64,
	) error
}

type EVMClient interface {
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
}

type Updater struct {
	logger           *slog.Logger
	l1Client         EVMClient
	l2Client         EVMClient
	winnerRegister   WinnerRegister
	evtMgr           events.EventManager
	l1BlockCache     *lru.Cache[uint64, map[string]int]
	l2BlockTimeCache *lru.Cache[uint64, uint64]
	metrics          *metrics
}

func NewUpdater(
	logger *slog.Logger,
	l1Client EVMClient,
	l2Client EVMClient,
	winnerRegister WinnerRegister,
	evtMgr events.EventManager,
) (*Updater, error) {
	l1BlockCache, err := lru.New[uint64, map[string]int](1024)
	if err != nil {
		return nil, fmt.Errorf("failed to create L1 block cache: %w", err)
	}
	l2BlockTimeCache, err := lru.New[uint64, uint64](1024)
	if err != nil {
		return nil, fmt.Errorf("failed to create L2 block time cache: %w", err)
	}
	return &Updater{
		logger:           logger,
		l1Client:         l1Client,
		l2Client:         l2Client,
		l1BlockCache:     l1BlockCache,
		l2BlockTimeCache: l2BlockTimeCache,
		winnerRegister:   winnerRegister,
		evtMgr:           evtMgr,
		metrics:          newMetrics(),
	}, nil
}

func (u *Updater) Metrics() []prometheus.Collector {
	return u.metrics.Collectors()
}

func (u *Updater) Start(ctx context.Context) <-chan struct{} {
	doneChan := make(chan struct{})

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return u.subscribeEncryptedCommitments(egCtx)
	})

	eg.Go(func() error {
		return u.subscribeOpenedCommitments(egCtx)
	})

	go func() {
		defer close(doneChan)
		if err := eg.Wait(); err != nil {
			u.logger.Error("failed to start updater", "error", err)
		}
	}()

	return doneChan
}

func (u *Updater) subscribeEncryptedCommitments(ctx context.Context) error {
	ev := events.NewEventHandler(
		"EncryptedCommitmentStored",
		func(update *preconf.PreconfcommitmentstoreEncryptedCommitmentStored) error {
			err := u.winnerRegister.AddEncryptedCommitment(
				ctx,
				update.CommitmentIndex[:],
				update.Commiter.Bytes(),
				update.CommitmentDigest[:],
				update.CommitmentSignature,
				update.BlockCommitedAt.Int64(),
			)
			if err != nil {
				u.logger.Error(
					"failed to add encrypted commitment",
					"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
					"error", err,
				)
				return err
			}
			u.metrics.EncryptedCommitmentsCount.Inc()
			return nil
		},
	)

	sub, err := u.evtMgr.Subscribe(ev)
	if err != nil {
		return fmt.Errorf("failed to subscribe to encrypted commitments: %w", err)
	}
	defer sub.Unsubscribe()

	select {
	case <-ctx.Done():
		return nil
	case err := <-sub.Err():
		return err
	}
}

func (u *Updater) subscribeOpenedCommitments(ctx context.Context) error {
	ev := events.NewEventHandler(
		"CommitmentStored",
		func(update *preconf.PreconfcommitmentstoreCommitmentStored) error {
			alreadySettled, err := u.winnerRegister.IsSettled(ctx, update.CommitmentIndex[:])
			if err != nil {
				u.logger.Error(
					"failed to check if commitment is settled",
					"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
					"error", err,
				)
				return err
			}
			if alreadySettled {
				// both bidders and providers could open commitments, so this could
				// be a duplicate event
				return nil
			}

			winner, err := u.winnerRegister.GetWinner(ctx, int64(update.BlockNumber))
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					log.Warn("winner not found", "blockNumber", update.BlockNumber)
					return nil
				}
				u.logger.Error(
					"failed to get winner",
					"blockNumber", update.BlockCommitedAt.Int64(),
					"error", err,
				)
				return err
			}

			if common.BytesToAddress(winner.Winner).Cmp(update.Commiter) != 0 {
				// The winner is not the committer of the commitment
				u.logger.Info(
					"winner is not the committer",
					"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
					"winner", winner.Winner,
					"committer", update.Commiter.Hex(),
				)
				return nil
			}

			txns, err := u.getL1Txns(ctx, update.BlockNumber)
			if err != nil {
				u.logger.Error(
					"failed to get L1 txns",
					"blockNumber", update.BlockNumber,
					"error", err,
				)
				return err
			}

			commitmentTxnHashes := strings.Split(update.TxnHash, ",")
			// Ensure Bundle is atomic and present in the block
			for i := 0; i < len(commitmentTxnHashes); i++ {
				posInBlock, found := txns[commitmentTxnHashes[i]]
				if !found || posInBlock != txns[commitmentTxnHashes[0]]+i {
					u.logger.Info(
						"bundle is not atomic",
						"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
						"txnHash", update.TxnHash,
					)
					// The committer did not include the transactions in the block
					// correctly, so this is a slash to be processed
					return u.addSettlement(
						ctx,
						update,
						settler.SettlementTypeSlash,
						0,
						winner.Window,
					)
				}
			}

			// Add the commitment as a reward
			// Compute the decay percentage
			l2BlockTime, err := u.getL2BlockTime(ctx, update.BlockCommitedAt.Uint64())
			if err != nil {
				u.logger.Error(
					"failed to get L2 block time",
					"blockNumber", update.BlockCommitedAt.Int64(),
					"error", err,
				)
				return err
			}

			decayPercentage := computeDecayPercentage(
				update.DecayStartTimeStamp,
				update.DecayEndTimeStamp,
				l2BlockTime,
			)

			return u.addSettlement(
				ctx,
				update,
				settler.SettlementTypeReward,
				decayPercentage,
				winner.Window,
			)
		},
	)

	sub, err := u.evtMgr.Subscribe(ev)
	if err != nil {
		return fmt.Errorf("failed to subscribe to opened commitments: %w", err)
	}
	defer sub.Unsubscribe()

	select {
	case <-ctx.Done():
		return nil
	case err := <-sub.Err():
		return err
	}
}

func (u *Updater) addSettlement(
	ctx context.Context,
	update *preconf.PreconfcommitmentstoreCommitmentStored,
	settlementType settler.SettlementType,
	decayPercentage int64,
	window int64,
) error {
	err := u.winnerRegister.AddSettlement(
		ctx,
		update.CommitmentIndex[:],
		update.TxnHash,
		int64(update.BlockNumber),
		update.Bid,
		update.Commiter.Bytes(),
		update.CommitmentHash[:],
		settlementType,
		decayPercentage,
		window,
	)
	if err != nil {
		u.logger.Error(
			"failed to add settlement",
			"commitmentIdx", common.Bytes2Hex(update.CommitmentIndex[:]),
			"error", err,
		)
		return err
	}

	u.metrics.CommitmentsCount.Inc()
	switch settlementType {
	case settler.SettlementTypeReward:
		u.metrics.RewardsCount.Inc()
	case settler.SettlementTypeSlash:
		u.metrics.SlashesCount.Inc()
	}

	return nil
}

func (u *Updater) getL1Txns(ctx context.Context, blockNum uint64) (map[string]int, error) {
	txns, ok := u.l1BlockCache.Get(blockNum)
	if ok {
		return txns, nil
	}

	blk, err := u.l1Client.BlockByNumber(ctx, big.NewInt(0).SetUint64(blockNum))
	if err != nil {
		return nil, fmt.Errorf("failed to get block by number: %w", err)
	}

	txnsInBlock := make(map[string]int)
	for posInBlock, tx := range blk.Transactions() {
		txnsInBlock[strings.TrimPrefix(tx.Hash().Hex(), "0x")] = posInBlock
	}
	_ = u.l1BlockCache.Add(blockNum, txnsInBlock)

	return txnsInBlock, nil
}

func (u *Updater) getL2BlockTime(ctx context.Context, blockNum uint64) (uint64, error) {
	time, ok := u.l2BlockTimeCache.Get(blockNum)
	if ok {
		return time, nil
	}

	blk, err := u.l2Client.BlockByNumber(ctx, big.NewInt(0).SetUint64(blockNum))
	if err != nil {
		return 0, fmt.Errorf("failed to get block by number: %w", err)
	}

	_ = u.l2BlockTimeCache.Add(blockNum, blk.Header().Time)

	return blk.Header().Time, nil
}

// computeDecayPercentage takes startTimestamp, endTimestamp, commitTimestamp and computes a linear decay percentage
// The computation does not care what format the timestamps are in, as long as they are consistent
// (e.g they could be unix or unixMili timestamps)
func computeDecayPercentage(startTimestamp, endTimestamp, commitTimestamp uint64) int64 {
	if startTimestamp >= endTimestamp || startTimestamp > commitTimestamp {
		return 0
	}

	// Calculate the total time in seconds
	totalTime := endTimestamp - startTimestamp
	// Calculate the time passed in seconds
	timePassed := commitTimestamp - startTimestamp
	// Calculate the decay percentage
	decayPercentage := float64(timePassed) / float64(totalTime)

	return int64(math.Round(decayPercentage * 100))
}
