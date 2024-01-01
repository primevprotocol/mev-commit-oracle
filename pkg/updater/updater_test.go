package updater_test

import (
	"context"
	"errors"
	"hash"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/primevprotocol/mev-oracle/pkg/preconf"
	"github.com/primevprotocol/mev-oracle/pkg/updater"
	"golang.org/x/crypto/sha3"
)

func getIdxBytes(idx int64) [32]byte {
	var idxBytes [32]byte
	big.NewInt(idx).FillBytes(idxBytes[:])
	return idxBytes
}

type testHasher struct {
	hasher hash.Hash
}

// NewHasher returns a new testHasher instance.
func NewHasher() *testHasher {
	return &testHasher{hasher: sha3.NewLegacyKeccak256()}
}

// Reset resets the hash state.
func (h *testHasher) Reset() {
	h.hasher.Reset()
}

// Update updates the hash state with the given key and value.
func (h *testHasher) Update(key, val []byte) error {
	h.hasher.Write(key)
	h.hasher.Write(val)
	return nil
}

// Hash returns the hash value.
func (h *testHasher) Hash() common.Hash {
	return common.BytesToHash(h.hasher.Sum(nil))
}

func TestUpdater(t *testing.T) {
	t.Parallel()

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	builderAddr := common.HexToAddress("0xabcd")

	signer := types.NewLondonSigner(big.NewInt(5))
	var txns []*types.Transaction
	for i := 0; i < 10; i++ {
		txns = append(txns, types.MustSignNewTx(key, signer, &types.DynamicFeeTx{
			Nonce:     uint64(i + 1),
			Gas:       1000000,
			Value:     big.NewInt(1),
			GasTipCap: big.NewInt(500),
			GasFeeCap: big.NewInt(500),
		}))
	}

	commitments := make(map[string]preconf.PreConfCommitmentStorePreConfCommitment)
	for i, txn := range txns {
		idxBytes := getIdxBytes(int64(i))
		commitments[string(idxBytes[:])] = preconf.PreConfCommitmentStorePreConfCommitment{
			Commiter: builderAddr,
			TxnHash:  txn.Hash().Hex(),
		}
	}

	testWinnerRegister := &testWinnerRegister{
		winners:     make(chan updater.BlockWinner),
		settlements: make(chan testSettlement),
		done:        make(chan int64, 1),
	}

	testL1Client := &testL1Client{
		blockNum: 5,
		block:    types.NewBlock(&types.Header{}, txns, nil, nil, NewHasher()),
	}

	testOracle := &testOracle{
		builder:     "test",
		builderAddr: builderAddr,
	}

	testPreconf := &testPreconf{
		blockNum:    5,
		commitments: commitments,
	}

	updtr := updater.NewUpdater(
		testL1Client,
		testWinnerRegister,
		testOracle,
		testPreconf,
	)

	ctx, cancel := context.WithCancel(context.Background())
	done := updtr.Start(ctx)

	testWinnerRegister.winners <- updater.BlockWinner{
		BlockNumber: 5,
		Winner:      "test",
	}

	count := 0
	for {
		if count == 10 {
			break
		}
		select {
		case settlement := <-testWinnerRegister.settlements:
			if settlement.blockNum != 5 {
				t.Fatal("wrong block number")
			}
			if settlement.builder != "test" {
				t.Fatal("wrong builder")
			}
			if settlement.isSlash {
				t.Fatal("should not be slash")
			}
			count++
		}
	}

	select {
	case <-testWinnerRegister.done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}

	cancel()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}
}

type testSettlement struct {
	commitmentIdx []byte
	txHash        string
	blockNum      int64
	builder       string
	isSlash       bool
}

type testWinnerRegister struct {
	winners     chan updater.BlockWinner
	settlements chan testSettlement
	done        chan int64
}

func (t *testWinnerRegister) SubscribeWinners(ctx context.Context) <-chan updater.BlockWinner {
	return t.winners
}

func (t *testWinnerRegister) UpdateComplete(ctx context.Context, blockNum int64) error {
	t.done <- blockNum
	return nil
}

func (t *testWinnerRegister) AddSettlement(
	ctx context.Context,
	commitmentIdx []byte,
	txHash string,
	blockNum int64,
	builder string,
	isSlash bool,
) error {
	t.settlements <- testSettlement{
		commitmentIdx: commitmentIdx,
		txHash:        txHash,
		blockNum:      blockNum,
		builder:       builder,
		isSlash:       isSlash,
	}
	return nil
}

type testL1Client struct {
	blockNum int64
	block    *types.Block
}

func (t *testL1Client) BlockByNumber(ctx context.Context, blkNum *big.Int) (*types.Block, error) {
	if blkNum.Int64() == t.blockNum {
		return t.block, nil
	}
	return nil, errors.New("block not found")
}

type testOracle struct {
	builder     string
	builderAddr common.Address
}

func (t *testOracle) GetBuilder(builder string) (common.Address, error) {
	if builder == t.builder {
		return t.builderAddr, nil
	}
	return common.Address{}, errors.New("builder not found")
}

type testPreconf struct {
	blockNum    int64
	commitments map[string]preconf.PreConfCommitmentStorePreConfCommitment
}

func (t *testPreconf) GetCommitmentsByBlockNumber(blockNum *big.Int) ([][32]byte, error) {
	if blockNum.Int64() == t.blockNum {
		var commitments [][32]byte
		for idx := range t.commitments {
			cIdx := [32]byte{}
			copy(cIdx[:], []byte(idx))
			commitments = append(commitments, cIdx)
		}
		return commitments, nil
	}

	return nil, errors.New("block not found")
}

func (t *testPreconf) GetCommitment(
	commitmentIdx [32]byte,
) (preconf.PreConfCommitmentStorePreConfCommitment, error) {
	if commitment, ok := t.commitments[string(commitmentIdx[:])]; ok {
		return commitment, nil
	}
	return preconf.PreConfCommitmentStorePreConfCommitment{}, errors.New("commitment not found")
}
