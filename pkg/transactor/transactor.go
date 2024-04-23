package transactor

import (
	"context"
	"log/slog"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/prometheus/client_golang/prometheus"
)

var allowedPendingTxnCount = 128

type TxnStore interface {
	LastNonce() (int64, error)
	SentTxn(nonce uint64, txHash common.Hash) error
	MarkSettlementComplete(ctx context.Context, nonce uint64) (int, error)
}

type EVMClient interface {
	bind.ContractTransactor
	BlockNumber(ctx context.Context) (uint64, error)
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
}

type Transactor struct {
	EVMClient
	store              TxnStore
	owner              common.Address
	nonceLock          sync.Mutex
	nonce              uint64
	qCond              *sync.Cond
	queue              []uint64
	confirmedNonceLock sync.Mutex
	confirmedNonce     uint64
	confirmedCond      *sync.Cond
	logger             *slog.Logger
	metrics            *metrics
}

func NewContractTransactor(
	ctx context.Context,
	client EVMClient,
	store TxnStore,
	owner common.Address,
	logger *slog.Logger,
) (*Transactor, error) {
	nonce, err := store.LastNonce()
	if err != nil {
		return nil, err
	}

	t := &Transactor{
		EVMClient: client,
		store:     store,
		owner:     owner,
		logger:    logger,
		nonce:     uint64(nonce),
		queue:     make([]uint64, 0, allowedPendingTxnCount),
		metrics:   newMetrics(),
	}

	t.qCond = sync.NewCond(&t.nonceLock)
	t.confirmedCond = sync.NewCond(&t.confirmedNonceLock)

	go t.txStatusUpdater(ctx)

	return t, nil
}

func (t *Transactor) Metrics() []prometheus.Collector {
	return t.metrics.Collectors()
}

func (t *Transactor) txStatusUpdater(ctx context.Context) {
	queryTicker := time.NewTicker(500 * time.Millisecond)
	defer queryTicker.Stop()

	lastBlock := uint64(0)
	for {
		select {
		case <-ctx.Done():
			return
		case <-queryTicker.C:
		}

		currentBlock, err := t.EVMClient.BlockNumber(ctx)
		if err != nil {
			t.logger.Error("failed to get block number", "error", err)
			continue
		}

		if currentBlock <= lastBlock {
			continue
		}

		lastNonce, err := t.EVMClient.NonceAt(
			ctx,
			t.owner,
			new(big.Int).SetUint64(currentBlock),
		)
		if err != nil {
			t.logger.Error("failed to get nonce", "error", err)
			continue
		}

		t.confirmedNonceLock.Lock()
		t.confirmedNonce = lastNonce
		t.confirmedCond.Broadcast()
		t.confirmedNonceLock.Unlock()

		t.metrics.LastConfirmedNonce.Set(float64(lastNonce))

		count, err := t.store.MarkSettlementComplete(ctx, lastNonce)
		if err != nil {
			t.logger.Error("failed to mark settlement complete", "error", err)
			continue
		}

		if count > 0 {
			t.logger.Info("marked settlement complete", "count", count)
		}

		lastBlock = currentBlock
	}
}

func (t *Transactor) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	t.nonceLock.Lock()
	defer t.nonceLock.Unlock()

	nonce, err := t.EVMClient.PendingNonceAt(ctx, account)
	if err != nil {
		return 0, err
	}

	if t.nonce < nonce {
		t.nonce = nonce
	}

	nonceToReturn := t.nonce
	t.nonce++
	t.queue = append(t.queue, nonceToReturn)

	t.metrics.LastUsedNonce.Set(float64(nonceToReturn))

	return nonceToReturn, nil
}

func (t *Transactor) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	// Wait if we have too many pending transactions
	t.confirmedNonceLock.Lock()
	for tx.Nonce()-t.confirmedNonce > uint64(allowedPendingTxnCount) {
		t.confirmedCond.Wait()
	}
	t.confirmedNonceLock.Unlock()

	// Wait for the transaction to be next in the queue
	t.nonceLock.Lock()
	for t.queue[0] != tx.Nonce() {
		t.qCond.Wait()
	}

	defer t.nonceLock.Unlock()

	err := t.EVMClient.SendTransaction(ctx, tx)
	if err != nil {
		return err
	}

	t.queue = t.queue[1:]
	t.qCond.Broadcast()

	t.metrics.LastSentNonce.Set(float64(tx.Nonce()))

	return t.store.SentTxn(tx.Nonce(), tx.Hash())
}
