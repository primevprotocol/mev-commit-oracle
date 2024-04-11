package settler_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	blocktracker "github.com/primevprotocol/contracts-abi/clients/BlockTracker"
	"github.com/primevprotocol/mev-oracle/pkg/events"
	"github.com/primevprotocol/mev-oracle/pkg/settler"
)

type testRegister struct {
	windowSettlements    map[int64][]settler.Settlement
	settlementsInitiated atomic.Int32
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

func TestSettler(t *testing.T) {
	t.Parallel()

	settlements := make(map[int64][]settler.Settlement)

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
	}

	orcl := &testOracle{
		commitments: make(chan processedCommitment, 10),
	}
	reg := &testRegister{
		windowSettlements: settlements,
	}

	btABI, err := abi.JSON(strings.NewReader(blocktracker.BlocktrackerABI))
	if err != nil {
		t.Fatal(err)
	}

	em := &testEventManager{
		btABI:     &btABI,
		windowSub: make(chan struct{}),
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
	}

	cancel()
	<-done
}

type testEventManager struct {
	btABI         *abi.ABI
	windowHandler events.EventHandler
	windowSub     chan struct{}
	sub           *testSub
}

type testSub struct {
	errC chan error
}

func (t *testSub) Unsubscribe() {}

func (t *testSub) Err() <-chan error {
	return t.errC
}

func (t *testEventManager) Subscribe(evt events.EventHandler) (events.Subscription, error) {
	if evt.EventName() == "NewWindow" {
		evt.SetTopicAndContract(t.btABI.Events["NewWindow"].ID, t.btABI)
		t.windowHandler = evt
		close(t.windowSub)
		return t.sub, nil
	}

	return nil, fmt.Errorf("event %s not found", evt.EventName())
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

	return t.windowHandler.Handle(testLog)
}
