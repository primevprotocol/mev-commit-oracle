package l1Listener_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	blocktracker "github.com/primevprotocol/contracts-abi/clients/BlockTracker"
	"github.com/primevprotocol/mev-oracle/pkg/events"
	"github.com/primevprotocol/mev-oracle/pkg/l1Listener"
)

func TestL1Listener(t *testing.T) {
	t.Parallel()

	reg := &testRegister{
		winners: make(chan winnerObj),
	}
	ethClient := &testEthClient{
		headers: make(map[uint64]*types.Header),
		errC:    make(chan error, 1),
	}
	oracle := &testOracle{
		builderMap: make(map[string]common.Address),
	}
	btABI, err := abi.JSON(strings.NewReader(blocktracker.BlocktrackerABI))
	if err != nil {
		t.Fatal(err)
	}
	eventManager := &testEventManager{
		btABI: &btABI,
		sub:   &testSub{errC: make(chan error)},
	}
	rec := &testRecorder{
		updates: make(chan l1Update),
	}

	l := l1Listener.NewL1Listener(
		slog.New(slog.NewTextHandler(io.Discard, nil)),
		ethClient,
		reg,
		oracle,
		eventManager,
		rec,
	)
	ctx, cancel := context.WithCancel(context.Background())

	cl := l1Listener.SetCheckInterval(100 * time.Millisecond)
	t.Cleanup(cl)

	done := l.Start(ctx)

	for i := 1; i < 10; i++ {
		ethClient.AddHeader(uint64(i), &types.Header{
			Number: big.NewInt(int64(i)),
			Extra:  []byte(fmt.Sprintf("b%d", i)),
		})
		addr := common.HexToAddress(fmt.Sprintf("0x%d", i))
		oracle.AddBuilder(fmt.Sprintf("b%d", i), addr)

		select {
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for winner", i)
		case update := <-rec.updates:
			if update.blockNum.Int64() != int64(i) {
				t.Fatal("wrong block number")
			}
			if update.winner.Cmp(addr) != 0 {
				t.Fatal("wrong winner")
			}
		}
	}

	// no winner
	ethClient.AddHeader(10, &types.Header{
		Number: big.NewInt(10),
	})

	// error registering winner, ensure it is retried
	ethClient.errC <- errors.New("dummy error")
	ethClient.AddHeader(11, &types.Header{
		Number: big.NewInt(11),
		Extra:  []byte("b11"),
	})
	addr := common.HexToAddress("0x11")
	oracle.AddBuilder("b11", addr)

	// ensure no winner is sent for the previous block
	select {
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for winner")
	case update := <-rec.updates:
		if update.blockNum.Int64() != 11 {
			t.Fatal("wrong block number")
		}
		if update.winner.Cmp(addr) != 0 {
			t.Fatal("wrong winner")
		}
	}

	for i := 1; i < 10; i++ {
		addr := common.HexToAddress(fmt.Sprintf("0x%d", i))
		go func() {
			err := eventManager.publish(
				big.NewInt(int64(i)),
				addr,
				big.NewInt(int64(i)),
			)
			if err != nil {
				t.Error(err)
			}
		}()

		select {
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for winner", i)
		case winner := <-reg.winners:
			if winner.blockNum != int64(i) {
				t.Fatal("wrong block number")
			}
			if !bytes.Equal(winner.winner, addr.Bytes()) {
				t.Fatal("wrong winner")
			}
			if winner.window != int64(i) {
				t.Fatal("wrong window")
			}
		}
	}

	cancel()
	select {
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for done")
	case <-done:
	}
}

type winnerObj struct {
	blockNum int64
	winner   []byte
	window   int64
}

type testRegister struct {
	winners chan winnerObj
}

func (t *testRegister) RegisterWinner(_ context.Context, blockNum int64, winner []byte, window int64) error {
	t.winners <- winnerObj{blockNum: blockNum, winner: winner, window: window}
	return nil
}

type testEthClient struct {
	mu      sync.Mutex
	headers map[uint64]*types.Header
	errC    chan error
}

func (t *testEthClient) AddHeader(blockNum uint64, hdr *types.Header) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.headers[blockNum] = hdr
}

func (t *testEthClient) BlockNumber(_ context.Context) (uint64, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(t.headers) == 0 {
		return 0, nil
	}
	blks := make([]uint64, len(t.headers))
	for k := range t.headers {
		blks = append(blks, k)
	}

	sort.Slice(blks, func(i, j int) bool {
		return blks[i] < blks[j]
	})

	return blks[len(blks)-1], nil
}

func (t *testEthClient) HeaderByNumber(_ context.Context, number *big.Int) (*types.Header, error) {
	select {
	case err := <-t.errC:
		return nil, err
	default:
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	hdr, ok := t.headers[number.Uint64()]
	if !ok {
		return nil, errors.New("header not found")
	}
	return hdr, nil
}

type testOracle struct {
	mu         sync.Mutex
	builderMap map[string]common.Address
}

func (t *testOracle) AddBuilder(builder string, addr common.Address) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.builderMap[builder] = addr
}

func (t *testOracle) GetBuilder(builder string) (common.Address, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	addr, ok := t.builderMap[builder]
	if !ok {
		return common.Address{}, errors.New("builder not found")
	}
	return addr, nil
}

type testEventManager struct {
	mu      sync.Mutex
	btABI   *abi.ABI
	handler events.EventHandler
	sub     *testSub
}

type testSub struct {
	errC chan error
}

func (t *testSub) Unsubscribe() {}

func (t *testSub) Err() <-chan error {
	return t.errC
}

func (t *testEventManager) Subscribe(evt events.EventHandler) (events.Subscription, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if evt.EventName() != "NewL1Block" {
		return nil, errors.New("invalid event")
	}
	evt.SetTopicAndContract(t.btABI.Events["NewL1Block"].ID, t.btABI)
	t.handler = evt
	return t.sub, nil
}

func (t *testEventManager) publish(blockNum *big.Int, winner common.Address, window *big.Int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	eventSignature := []byte("NewL1Block(uint256,address,uint256)")
	hashEventSignature := crypto.Keccak256Hash(eventSignature)

	blockNumber := common.BigToHash(blockNum)
	winnerHash := common.HexToHash(winner.Hex())
	windowNumber := common.BigToHash(window)

	// Creating a Log object
	testLog := types.Log{
		Topics: []common.Hash{
			hashEventSignature, // The first topic is the hash of the event signature
			blockNumber,        // The next topics are the indexed event parameters
			winnerHash,
			windowNumber,
		},
		// Since there are no non-indexed parameters, Data is empty
		Data: []byte{},
	}

	return t.handler.Handle(testLog)
}

type l1Update struct {
	blockNum *big.Int
	winner   common.Address
}

type testRecorder struct {
	updates chan l1Update
}

func (t *testRecorder) RecordL1Block(blockNum *big.Int, addr common.Address) (*types.Transaction, error) {
	t.updates <- l1Update{blockNum: blockNum, winner: addr}
	return types.NewTransaction(0, addr, nil, 0, nil, nil), nil
}
