package transactor_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/primevprotocol/mev-oracle/pkg/transactor"
	"golang.org/x/sync/errgroup"
)

func TestTransactor(t *testing.T) {
	t.Parallel()

	store := &testStore{}
	evmClient := &testEVMClient{}
	owner := common.HexToAddress("0x1234")

	cleanup := transactor.SetAllowedPendingTxnCount(5)
	t.Cleanup(cleanup)

	ctx := context.Background()

	txtor, err := transactor.NewContractTransactor(
		ctx,
		evmClient,
		store,
		owner,
		slog.New(slog.NewTextHandler(io.Discard, nil)),
	)
	if err != nil {
		t.Fatal(err)
	}

	for i := 1; i <= 5; i++ {
		nonce, err := txtor.PendingNonceAt(ctx, owner)
		if err != nil {
			t.Fatal(err)
		}

		if nonce != uint64(i) {
			t.Fatalf("expected nonce %d, got %d", i, nonce)
		}
	}

	eg := errgroup.Group{}

	for i := 5; i >= 1; i-- {
		tx := types.NewTransaction(uint64(i), owner, big.NewInt(0), 21000, big.NewInt(1), nil)
		eg.Go(func() error {
			if err := txtor.SendTransaction(ctx, tx); err != nil {
				return err
			}
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}

	sentTxnNonces := store.GetSentTxns()
	if len(sentTxnNonces) != 5 {
		t.Fatalf("expected 5 sent transactions, got %d", len(sentTxnNonces))
	}

	for i, nonce := range sentTxnNonces {
		if nonce != uint64(i+1) {
			t.Fatalf("expected nonce %d, got %d", i+1, nonce)
		}
	}

	for i := 6; i <= 10; i++ {
		nonce, err := txtor.PendingNonceAt(ctx, owner)
		if err != nil {
			t.Fatal(err)
		}

		if nonce != uint64(i) {
			t.Fatalf("expected nonce %d, got %d", i, nonce)
		}

		tx := types.NewTransaction(uint64(i), owner, big.NewInt(0), 21000, big.NewInt(1), nil)
		eg.Go(func() error {
			if err := txtor.SendTransaction(ctx, tx); err != nil {
				return err
			}
			return nil
		})
	}

	// this sleep is to ensure that we wait enough time to ensure the transactions
	// are waiting in the queue
	time.Sleep(1 * time.Second)

	sentTxnNonces = store.GetSentTxns()
	if len(sentTxnNonces) != 5 {
		t.Fatalf("expected 5 sent transactions, got %d", len(sentTxnNonces))
	}

	evmClient.SetBlockNumber(1)
	evmClient.SetNonce(5)

	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}

	sentTxnNonces = store.GetSentTxns()
	if len(sentTxnNonces) != 10 {
		t.Fatalf("expected 10 sent transactions, got %d", len(sentTxnNonces))
	}

	for i, nonce := range sentTxnNonces {
		if nonce != uint64(i+1) {
			t.Fatalf("expected nonce %d, got %d", i+1, nonce)
		}
	}

	confirmedNonce := store.GetConfirmedNonce()
	if confirmedNonce != 5 {
		t.Fatalf("expected confirmed nonce 5, got %d", confirmedNonce)
	}
}

type testStore struct {
	mu             sync.Mutex
	confirmedNonce uint64
	sentTxns       []uint64
}

func (t *testStore) LastNonce() (uint64, error) {
	return 0, nil
}

func (t *testStore) SentTxn(nonce uint64, txHash common.Hash) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.sentTxns = append(t.sentTxns, nonce)
	return nil
}

func (t *testStore) MarkSettlementComplete(ctx context.Context, nonce uint64) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.confirmedNonce = nonce
	return 0, nil
}

func (t *testStore) GetConfirmedNonce() uint64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.confirmedNonce
}

func (t *testStore) GetSentTxns() []uint64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.sentTxns
}

type testEVMClient struct {
	bind.ContractTransactor
	mu           sync.Mutex
	blockNumber  uint64
	nonce        uint64
	currentNonce uint64
}

func (t *testEVMClient) BlockNumber(ctx context.Context) (uint64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.blockNumber, nil
}

func (t *testEVMClient) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.nonce, nil
}

func (t *testEVMClient) SetBlockNumber(blockNumber uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.blockNumber = blockNumber
}

func (t *testEVMClient) SetNonce(nonce uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.nonce = nonce
}

func (t *testEVMClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	return 1, nil
}

func (t *testEVMClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.currentNonce+1 != tx.Nonce() {
		return fmt.Errorf("expected nonce %d, got %d", t.currentNonce+1, tx.Nonce())
	}

	t.currentNonce = tx.Nonce()

	return nil
}
