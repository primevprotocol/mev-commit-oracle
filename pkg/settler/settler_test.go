package settler_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	bidderregistry "github.com/primevprotocol/contracts-abi/clients/BidderRegistry"
	blocktracker "github.com/primevprotocol/contracts-abi/clients/BlockTracker"
	"github.com/primevprotocol/mev-oracle/pkg/events"
	"github.com/primevprotocol/mev-oracle/pkg/settler"
)

type testRegister struct {
	windowSettlements    map[int64][]settler.Settlement
	windowReturns        map[int64][]settler.Return
	settlementsInitiated atomic.Int32
	returnsInitiated     atomic.Int32
	bidderChan           chan registeredBidder
}

func (t *testRegister) SubscribeSettlements(ctx context.Context, window int64) <-chan settler.Settlement {
	sc := make(chan settler.Settlement)
	go func() {
		defer close(sc)
		settlements := t.windowSettlements[window]
		for _, s := range settlements {
			select {
			case <-ctx.Done():
				return
			case sc <- s:
			}
		}
	}()

	return sc
}

func (t *testRegister) SubscribeReturns(ctx context.Context, _ int, window int64) <-chan settler.Return {
	rc := make(chan settler.Return)
	go func() {
		defer close(rc)
		returns := t.windowReturns[window]
		for _, r := range returns {
			select {
			case <-ctx.Done():
				return
			case rc <- r:
			}
		}
	}()

	return rc
}

type registeredBidder struct {
	bidder []byte
	amount int64
	window int64
}

func (t *testRegister) BidderRegistered(ctx context.Context, bidder []byte, window int64, amount int64) error {
	t.bidderChan <- registeredBidder{
		bidder: bidder,
		amount: amount,
		window: window,
	}
	return nil
}

func (t *testRegister) ReturnInitiated(
	ctx context.Context,
	window int64,
	bidders [][]byte,
	txHash common.Hash,
	nonce uint64,
) error {
	t.returnsInitiated.Add(1)
	return nil
}

func (t *testRegister) SettlementInitiated(
	ctx context.Context,
	commitmentIdx []byte,
	txHash common.Hash,
	nonce uint64,
) error {
	t.settlementsInitiated.Add(1)
	return nil
}

type processedCommitment struct {
	commitmentIdx [32]byte
	blockNum      *big.Int
	builder       string
	isSlash       bool
	residualDecay *big.Int
}

type testOracle struct {
	commitments chan processedCommitment
	returns     chan common.Address
}

func (t *testOracle) ProcessBuilderCommitmentForBlockNumber(
	commitmentIdx [32]byte,
	blockNum *big.Int,
	builder string,
	isSlash bool,
	residualDecay *big.Int,
) (*types.Transaction, error) {
	t.commitments <- processedCommitment{
		commitmentIdx: commitmentIdx,
		blockNum:      blockNum,
		builder:       builder,
		isSlash:       isSlash,
		residualDecay: residualDecay,
	}
	return types.NewTransaction(0, common.Address{}, nil, 0, nil, nil), nil
}

func (t *testOracle) UnlockFunds(_ *big.Int, bidders []common.Address) (*types.Transaction, error) {
	for _, bidder := range bidders {
		t.returns <- bidder
	}
	return types.NewTransaction(0, common.Address{}, nil, 0, nil, nil), nil
}

func TestSettler(t *testing.T) {
	t.Parallel()

	settlements := make(map[int64][]settler.Settlement)
	returns := make(map[int64][]settler.Return)

	for i := 0; i < 10; i++ {
		var sType settler.SettlementType
		if i%2 == 0 {
			sType = settler.SettlementTypeReward
		} else {
			sType = settler.SettlementTypeSlash
		}
		settlements[100] = append(settlements[100], settler.Settlement{
			CommitmentIdx: big.NewInt(int64(i + 1)).Bytes(),
			TxHash:        "0x1234",
			BlockNum:      100,
			Builder:       common.HexToAddress("0x1234").Bytes(),
			Amount:        1000,
			BidID:         common.HexToHash(fmt.Sprintf("0x%02d", i)).Bytes(),
			Type:          sType,
		})

		returns[99] = append(returns[99], settler.Return{
			Window:  99,
			Bidders: [][]byte{common.HexToAddress("0x1234").Bytes()},
		})
	}

	orcl := &testOracle{
		commitments: make(chan processedCommitment, 10),
		returns:     make(chan common.Address, 10),
	}
	reg := &testRegister{
		windowSettlements: settlements,
		windowReturns:     returns,
		bidderChan:        make(chan registeredBidder, 1),
	}

	btABI, err := abi.JSON(strings.NewReader(blocktracker.BlocktrackerABI))
	if err != nil {
		t.Fatal(err)
	}

	bidderABI, err := abi.JSON(strings.NewReader(bidderregistry.BidderregistryABI))
	if err != nil {
		t.Fatal(err)
	}

	em := &testEventManager{
		btABI:     &btABI,
		bidderABI: &bidderABI,
		windowSub: make(chan struct{}),
		bidderSub: make(chan struct{}),
		sub:       &testSub{errC: make(chan error)},
	}

	s := settler.NewSettler(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		orcl,
		reg,
		em,
	)

	ctx, cancel := context.WithCancel(context.Background())
	done := s.Start(ctx)

	<-em.windowSub
	<-em.bidderSub

	for i := 0; i < 10; i++ {
		// Test that the settler is able to process a bidder registration
		b := bidderregistry.BidderregistryBidderRegistered{
			Bidder:        common.HexToAddress(fmt.Sprintf("0x%02d", i)),
			PrepaidAmount: big.NewInt(1000),
			WindowNumber:  big.NewInt(99),
		}
		if err := em.publishBidderRegistered(b); err != nil {
			t.Fatal(err)
		}

		select {
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for bidder registration")
		case r := <-reg.bidderChan:
			if r.amount != 1000 {
				t.Fatalf("expected amount 1000, got %d", r.amount)
			}
			if r.window != 99 {
				t.Fatalf("expected window 99, got %d", r.window)
			}
		}
	}

	// Test that the settler is able to process a window
	w := blocktracker.BlocktrackerNewWindow{
		Window: big.NewInt(102),
	}
	if err := em.publishNewWindow(w); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		select {
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for commitment")
		case s := <-orcl.commitments:
			if s.blockNum.Cmp(big.NewInt(100)) != 0 {
				t.Fatalf("expected block number 100, got %d", s.blockNum)
			}
			if common.HexToAddress(s.builder).Cmp(common.HexToAddress("0x1234")) != 0 {
				t.Fatalf(
					"expected builder %s, got %s",
					common.HexToAddress("0x1234"),
					common.HexToAddress(s.builder),
				)
			}
			if s.isSlash && i%2 == 0 {
				t.Fatalf("expected slash, got reward")
			}
			if !s.isSlash && i%2 != 0 {
				t.Fatalf("expected reward, got slash")
			}
		}

		select {
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for return")
		case r := <-orcl.returns:
			if r.Hex() != common.HexToAddress("0x1234").Hex() {
				t.Fatalf("expected bidder 0x1234, got %s", r.Hex())
			}
		}
	}

	cancel()
	<-done
}

type testEventManager struct {
	btABI          *abi.ABI
	bidderABI      *abi.ABI
	mu             sync.Mutex
	windowHandlers []events.EventHandler
	windowSub      chan struct{}
	bidderHandler  events.EventHandler
	bidderSub      chan struct{}
	sub            *testSub
}

type testSub struct {
	errC chan error
}

func (t *testSub) Unsubscribe() {}

func (t *testSub) Err() <-chan error {
	return t.errC
}

func (t *testEventManager) Subscribe(evt events.EventHandler) (events.Subscription, error) {
	switch evt.EventName() {
	case "NewWindow":
		t.mu.Lock()
		evt.SetTopicAndContract(t.btABI.Events["NewWindow"].ID, t.btABI)
		t.windowHandlers = append(t.windowHandlers, evt)
		t.mu.Unlock()
		if len(t.windowHandlers) == 2 {
			close(t.windowSub)
		}
	case "BidderRegistered":
		evt.SetTopicAndContract(t.bidderABI.Events["BidderRegistered"].ID, t.bidderABI)
		t.bidderHandler = evt
		close(t.bidderSub)
	default:
		return nil, fmt.Errorf("event %s not found", evt.EventName())
	}

	return t.sub, nil
}

func (t *testEventManager) publishNewWindow(w blocktracker.BlocktrackerNewWindow) error {
	event := t.btABI.Events["NewWindow"]

	// Creating a Log object
	testLog := types.Log{
		Topics: []common.Hash{
			event.ID,                   // The first topic is the hash of the event signature
			common.BigToHash(w.Window), // The next topics are the indexed event parameters
		},
		// Non-indexed parameters are stored in the Data field
		Data: nil,
	}

	for _, h := range t.windowHandlers {
		if err := h.Handle(testLog); err != nil {
			return err
		}
	}

	return nil
}

func (t *testEventManager) publishBidderRegistered(c bidderregistry.BidderregistryBidderRegistered) error {
	event := t.bidderABI.Events["BidderRegistered"]
	buf, err := event.Inputs.NonIndexed().Pack(
		c.PrepaidAmount,
		c.WindowNumber,
	)
	if err != nil {
		return err
	}

	bidder := common.HexToHash(c.Bidder.Hex())

	// Creating a Log object
	testLog := types.Log{
		Topics: []common.Hash{
			event.ID, // The first topic is the hash of the event signature
			bidder,   // The next topics are the indexed event parameters
		},
		Data: buf,
	}

	return t.bidderHandler.Handle(testLog)
}
