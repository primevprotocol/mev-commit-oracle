package events_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderregistry "github.com/primevprotocol/contracts-abi/clients/BidderRegistry"
	"github.com/primevprotocol/mev-oracle/pkg/events"
)

func TestEventHandler(t *testing.T) {
	t.Parallel()

	b := bidderregistry.BidderregistryBidderRegistered{
		Bidder:        common.HexToAddress("0xabcd"),
		PrepaidAmount: big.NewInt(1000),
		WindowNumber:  big.NewInt(99),
	}

	evtHdlr := events.NewEventHandler(
		"BidderRegistered",
		func(ev *bidderregistry.BidderregistryBidderRegistered) error {
			if ev.Bidder.Hex() != b.Bidder.Hex() {
				return fmt.Errorf("expected bidder %s, got %s", b.Bidder.Hex(), ev.Bidder.Hex())
			}
			if ev.PrepaidAmount.Cmp(b.PrepaidAmount) != 0 {
				return fmt.Errorf("expected prepaid amount %d, got %d", b.PrepaidAmount, ev.PrepaidAmount)
			}
			if ev.WindowNumber.Cmp(b.WindowNumber) != 0 {
				return fmt.Errorf("expected window number %d, got %d", b.WindowNumber, ev.WindowNumber)
			}
			return nil
		},
	)

	bidderABI, err := abi.JSON(strings.NewReader(bidderregistry.BidderregistryABI))
	if err != nil {
		t.Fatal(err)
	}

	event := bidderABI.Events["BidderRegistered"]

	evtHdlr.SetTopicAndContract(event.ID, &bidderABI)

	if evtHdlr.Topic().Cmp(event.ID) != 0 {
		t.Fatalf("expected topic %s, got %s", event.ID, evtHdlr.Topic())
	}

	if evtHdlr.EventName() != "BidderRegistered" {
		t.Fatalf("expected event name BidderRegistered, got %s", evtHdlr.EventName())
	}

	buf, err := event.Inputs.NonIndexed().Pack(
		b.PrepaidAmount,
		b.WindowNumber,
	)
	if err != nil {
		t.Fatal(err)
	}

	bidder := common.HexToHash(b.Bidder.Hex())

	// Creating a Log object
	testLog := types.Log{
		Topics: []common.Hash{
			event.ID, // The first topic is the hash of the event signature
			bidder,   // The next topics are the indexed event parameters
		},
		Data: buf,
	}

	if err := evtHdlr.Handle(testLog); err != nil {
		t.Fatal(err)
	}
}

func TestEventManager(t *testing.T) {
	t.Parallel()

	b := bidderregistry.BidderregistryBidderRegistered{
		Bidder:        common.HexToAddress("0xabcd"),
		PrepaidAmount: big.NewInt(1000),
		WindowNumber:  big.NewInt(99),
	}

	handlerTriggered := make(chan struct{})

	evtHdlr := events.NewEventHandler(
		"BidderRegistered",
		func(ev *bidderregistry.BidderregistryBidderRegistered) error {
			defer close(handlerTriggered)

			if ev.Bidder.Hex() != b.Bidder.Hex() {
				return fmt.Errorf("expected bidder %s, got %s", b.Bidder.Hex(), ev.Bidder.Hex())
			}
			if ev.PrepaidAmount.Cmp(b.PrepaidAmount) != 0 {
				return fmt.Errorf("expected prepaid amount %d, got %d", b.PrepaidAmount, ev.PrepaidAmount)
			}
			if ev.WindowNumber.Cmp(b.WindowNumber) != 0 {
				return fmt.Errorf("expected window number %d, got %d", b.WindowNumber, ev.WindowNumber)
			}
			return nil
		},
	)

	bidderABI, err := abi.JSON(strings.NewReader(bidderregistry.BidderregistryABI))
	if err != nil {
		t.Fatal(err)
	}

	evmClient := &testEVMClient{
		logsSub: make(chan struct{}),
		sub: &testSub{
			errC: make(chan error),
		},
	}

	store := &testStore{}

	contracts := map[common.Address]*abi.ABI{
		common.HexToAddress("0xabcd"): &bidderABI,
	}

	evtMgr := events.NewListener(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		evmClient,
		store,
		contracts,
	)

	ctx, cancel := context.WithCancel(context.Background())
	done := evtMgr.Start(ctx)

	sub, err := evtMgr.Subscribe(evtHdlr)
	if err != nil {
		t.Fatal(err)
	}

	defer sub.Unsubscribe()

	<-evmClient.logsSub

	data, err := bidderABI.Events["BidderRegistered"].Inputs.NonIndexed().Pack(
		b.PrepaidAmount,
		b.WindowNumber,
	)
	if err != nil {
		t.Fatal(err)
	}

	evmClient.logs <- types.Log{
		Topics: []common.Hash{
			bidderABI.Events["BidderRegistered"].ID,
			common.HexToHash(b.Bidder.Hex()),
		},
		Data:        data,
		BlockNumber: 1,
	}

	select {
	case <-handlerTriggered:
	case err := <-sub.Err():
		t.Fatal("handler not triggered", err)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for handler to be triggered")
	}

	if b, err := store.LastBlock(); err != nil || b != 1 {
		t.Fatalf("expected block number 1, got %d", store.blockNumber)
	}

	cancel()
	<-done
}

type testEVMClient struct {
	logs    chan<- types.Log
	logsSub chan struct{}
	sub     *testSub
}

type testSub struct {
	errC chan error
}

func (t *testSub) Unsubscribe() {}

func (t *testSub) Err() <-chan error {
	return t.errC
}

func (t *testEVMClient) SubscribeFilterLogs(
	ctx context.Context,
	q ethereum.FilterQuery,
	ch chan<- types.Log,
) (ethereum.Subscription, error) {
	defer close(t.logsSub)
	t.logs = ch
	return t.sub, nil
}

type testStore struct {
	mu          sync.Mutex
	blockNumber uint64
}

func (t *testStore) LastBlock() (uint64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.blockNumber, nil
}

func (t *testStore) SetLastBlock(blockNumber uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.blockNumber = blockNumber
	return nil
}
