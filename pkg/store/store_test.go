package store_test

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/primevprotocol/mev-oracle/pkg/settler"
	"github.com/primevprotocol/mev-oracle/pkg/store"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type blockWinner struct {
	BlockNumber int64
	Winner      []byte
	Window      int64
}

func TestStore(t *testing.T) {
	ctx := context.Background()

	// Define the PostgreSQL container request
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "user",
			"POSTGRES_PASSWORD": "password",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	// Start the PostgreSQL container
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %s", err)
	}
	defer func() {
		err := postgresContainer.Terminate(ctx)
		if err != nil {
			t.Errorf("Failed to terminate PostgreSQL container: %s", err)
		}
	}()

	// Retrieve the container's mapped port
	mappedPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get mapped port: %s", err)
	}
	// Construct the database connection string
	connStr := fmt.Sprintf("postgresql://user:password@localhost:%s/testdb?sslmode=disable", mappedPort.Port())

	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL container: %s", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping PostgreSQL container: %s", err)
	}

	winners := []blockWinner{
		{
			Window:      1,
			Winner:      common.HexToAddress("0x01").Bytes(),
			BlockNumber: 1,
		},
		{
			Window:      2,
			Winner:      common.HexToAddress("0x02").Bytes(),
			BlockNumber: 2,
		},
	}

	settlements := []settler.Settlement{
		{
			CommitmentIdx: []byte{1},
			TxHash:        common.HexToHash("0x01").String(),
			BlockNum:      1,
			Amount:        2000000,
			Builder:       winners[0].Winner,
			BidID:         common.HexToHash("0x01").Bytes(),
			Type:          settler.SettlementTypeReward,
		},
		{
			CommitmentIdx: []byte{2},
			TxHash:        common.HexToHash("0x02").String(),
			BlockNum:      1,
			Amount:        1000000,
			Builder:       winners[0].Winner,
			BidID:         common.HexToHash("0x02").Bytes(),
			Type:          settler.SettlementTypeSlash,
		},
		{
			CommitmentIdx: []byte{3},
			TxHash:        common.HexToHash("0x03").String(),
			BlockNum:      1,
			Amount:        1000000,
			Builder:       winners[1].Winner,
			BidID:         common.HexToHash("0x03").Bytes(),
			Type:          settler.SettlementTypeReturn,
		},
		{
			CommitmentIdx: []byte{4},
			TxHash:        common.HexToHash("0x04").String(),
			BlockNum:      2,
			Amount:        2000000,
			Builder:       winners[1].Winner,
			BidID:         common.HexToHash("0x04").Bytes(),
			Type:          settler.SettlementTypeReward,
		},
		{
			CommitmentIdx: []byte{5},
			TxHash:        common.HexToHash("0x05").String(),
			BlockNum:      2,
			Amount:        1000000,
			Builder:       winners[1].Winner,
			BidID:         common.HexToHash("0x05").Bytes(),
			Type:          settler.SettlementTypeSlash,
		},
		{
			CommitmentIdx: []byte{6},
			TxHash:        common.HexToHash("0x06").String(),
			BlockNum:      2,
			Amount:        1000000,
			Builder:       winners[0].Winner,
			BidID:         common.HexToHash("0x04").Bytes(),
			Type:          settler.SettlementTypeReturn,
		},
	}

	t.Run("NewStore", func(t *testing.T) {
		// Create the store and tables
		_, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}
	})

	t.Run("RegisterWinner", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		for _, winner := range winners {
			err = st.RegisterWinner(context.Background(), winner.BlockNumber, winner.Winner, winner.Window)
			if err != nil {
				t.Fatalf("Failed to register winner: %s", err)
			}
		}
	})

	t.Run("GetWinner", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		for _, winner := range winners {
			w, err := st.GetWinner(context.Background(), winner.BlockNumber)
			if err != nil {
				t.Fatalf("Failed to get winner: %s", err)
			}
			if diff := cmp.Diff(w.Winner, winner.Winner); diff != "" {
				t.Fatalf("Unexpected winner: (-want +have):\n%s", diff)
			}
		}
	})

	t.Run("AddEncryptedCommitment", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		for i, settlement := range settlements {
			blkNo := int64(i/3) + 1
			err = st.AddEncryptedCommitment(
				context.Background(),
				settlement.CommitmentIdx,
				settlement.Builder,
				[]byte("hash"),
				[]byte("signature"),
				blkNo,
			)
			if err != nil {
				t.Fatalf("Failed to add encrypted commitment: %s", err)
			}
		}
	})

	t.Run("AddSettlement", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		for i, settlement := range settlements {
			window := int64(i/3) + 1
			err = st.AddSettlement(
				context.Background(),
				settlement.CommitmentIdx,
				settlement.TxHash,
				settlement.BlockNum,
				settlement.Amount,
				settlement.Builder,
				settlement.BidID,
				settlement.Type,
				settlement.DecayPercentage,
				window,
			)
			if err != nil {
				t.Fatalf("Failed to add settlement: %s", err)
			}
		}
	})

	t.Run("SubscribeSettlements", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		ctx, cancel := context.WithCancel(context.Background())

		settlementChan := st.SubscribeSettlements(ctx, winners[0].Window)
		idx := 0
		for s := range settlementChan {
			if diff := cmp.Diff(s, settlements[idx]); diff != "" {
				t.Fatalf("Unexpected settlement: (-want +have):\n%s", diff)
			}
			idx++
		}

		idx++

		settlementChan2 := st.SubscribeSettlements(ctx, winners[1].Window)
		for s := range settlementChan2 {
			if diff := cmp.Diff(s, settlements[idx]); diff != "" {
				t.Fatalf("Unexpected settlement: (-want +have):\n%s", diff)
			}
			idx++
		}

		if idx != len(settlements)-1 {
			t.Fatalf("Expected %d settlements, got %d", len(settlements), idx)
		}

		cancel()
		sChan := st.SubscribeSettlements(ctx, winners[0].Window)
		_, ok := <-sChan
		if ok {
			t.Fatalf("Expected channel to be closed")
		}
	})

	t.Run("BidderRegistered", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		err = st.BidderRegistered(
			context.Background(),
			common.HexToAddress("0x1234").Bytes(),
			1,
			1000,
		)
		if err != nil {
			t.Fatalf("Failed to register bidder: %s", err)
		}
	})

	t.Run("SubscribeReturns", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		ctx, cancel := context.WithCancel(context.Background())

		returnChan := st.SubscribeReturns(ctx, 1, 1)

		for r := range returnChan {
			if !bytes.Equal(r.Bidders[0], common.HexToAddress("0x1234").Bytes()) {
				t.Fatalf(
					"Unexpected return settlement: want %s have %s\n",
					common.HexToAddress("0x1234"),
					common.BytesToAddress(r.Bidders[0]),
				)
			}
		}

		cancel()

		returnChan = st.SubscribeReturns(ctx, 1, 1)
		_, ok := <-returnChan
		if ok {
			t.Fatalf("Expected channel to be closed")
		}
	})

	t.Run("SettlementInitiated", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		for i := range []int{0, 1, 3, 4} {
			err = st.SentTxn(uint64(i+1), common.HexToHash(fmt.Sprintf("0x%02d", i)))
			if err != nil {
				t.Fatalf("Failed to mark txn sent: %s", err)
			}

			err = st.SettlementInitiated(
				context.Background(),
				settlements[i].CommitmentIdx,
				common.HexToHash(fmt.Sprintf("0x%02d", i)),
				uint64(i+1),
			)
			if err != nil {
				t.Fatalf("Failed to initiate settlement: %s", err)
			}
		}
	})

	t.Run("ReturnInitiated", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		err = st.SentTxn(5, common.HexToHash(fmt.Sprintf("0x%02d", 5)))
		if err != nil {
			t.Fatalf("Failed to mark txn sent: %s", err)
		}

		err = st.ReturnInitiated(
			context.Background(),
			1,
			[][]byte{common.HexToAddress("0x1234").Bytes()},
			common.HexToHash(fmt.Sprintf("0x%02d", 5)),
			uint64(5),
		)
		if err != nil {
			t.Fatalf("Failed to initiate settlement: %s", err)
		}
	})

	t.Run("LastNonce and PendingTxnCount", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		lastNonce, err := st.LastNonce()
		if err != nil {
			t.Fatalf("Failed to get last nonce: %s", err)
		}
		if lastNonce != 5 {
			t.Fatalf("Expected last nonce 5, got %d", lastNonce)
		}

		pendingTxnCount, err := st.PendingTxnCount()
		if err != nil {
			t.Fatalf("Failed to get pending txn count: %s", err)
		}
		if pendingTxnCount != 5 {
			t.Fatalf("Expected pending txn count 5, got %d", pendingTxnCount)
		}
	})

	t.Run("MarkSettlementComplete", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		count, err := st.MarkSettlementComplete(context.Background(), 6)
		if err != nil {
			t.Fatalf("Failed to mark settlement complete: %s", err)
		}
		if count != 5 {
			t.Fatalf("Expected count 5, got %d", count)
		}

		pendingTxnCount, err := st.PendingTxnCount()
		if err != nil {
			t.Fatalf("Failed to get pending txn count: %s", err)
		}
		if pendingTxnCount != 0 {
			t.Fatalf("Expected pending txn count 0, got %d", pendingTxnCount)
		}
	})
}
