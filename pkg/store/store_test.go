package store_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/primevprotocol/mev-oracle/pkg/settler"
	"github.com/primevprotocol/mev-oracle/pkg/store"
	"github.com/primevprotocol/mev-oracle/pkg/updater"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

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
	defer postgresContainer.Terminate(ctx)

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

	winners := []updater.BlockWinner{
		{
			BlockNumber: 1,
			Winner:      common.HexToAddress("0x01").String(),
		},
		{
			BlockNumber: 2,
			Winner:      common.HexToAddress("0x02").String(),
		},
	}

	settlements := []settler.Settlement{
		{
			CommitmentIdx: []byte{1},
			TxHash:        common.HexToHash("0x01").String(),
			BlockNum:      1,
			Amount:        2000000,
			Builder:       winners[0].Winner,
			Type:          settler.SettlementTypeReward,
		},
		{
			CommitmentIdx: []byte{2},
			TxHash:        common.HexToHash("0x02").String(),
			BlockNum:      1,
			Amount:        1000000,
			Builder:       winners[0].Winner,
			Type:          settler.SettlementTypeSlash,
		},
		{
			CommitmentIdx: []byte{3},
			TxHash:        common.HexToHash("0x03").String(),
			BlockNum:      1,
			Amount:        1000000,
			Builder:       winners[1].Winner,
			Type:          settler.SettlementTypeReturn,
		},
		{
			CommitmentIdx: []byte{4},
			TxHash:        common.HexToHash("0x04").String(),
			BlockNum:      2,
			Amount:        2000000,
			Builder:       winners[1].Winner,
			Type:          settler.SettlementTypeReward,
		},
		{
			CommitmentIdx: []byte{5},
			TxHash:        common.HexToHash("0x05").String(),
			BlockNum:      2,
			Amount:        1000000,
			Builder:       winners[1].Winner,
			Type:          settler.SettlementTypeSlash,
		},
		{
			CommitmentIdx: []byte{6},
			TxHash:        common.HexToHash("0x06").String(),
			BlockNum:      2,
			Amount:        1000000,
			Builder:       winners[0].Winner,
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
			err = st.RegisterWinner(context.Background(), winner.BlockNumber, winner.Winner)
			if err != nil {
				t.Fatalf("Failed to register winner: %s", err)
			}
		}
	})

	t.Run("SubscribeWinners", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		ctx, cancel := context.WithCancel(context.Background())

		// Subscribe to winners
		winnerChan := st.SubscribeWinners(ctx)
		if err != nil {
			t.Fatalf("Failed to subscribe to winners: %s", err)
		}

		for i := 0; i < 2; i++ {
			winner := <-winnerChan
			if winner.BlockNumber != winners[i].BlockNumber {
				t.Fatalf("Expected block number %d, got %d", winners[i].BlockNumber, winner.BlockNumber)
			}
			if winner.Winner != winners[i].Winner {
				t.Fatalf("Expected builder address %s, got %s", winners[i].Winner, winner.Winner)
			}
		}

		cancel()

		winner, ok := <-winnerChan
		if ok {
			t.Fatalf("Expected channel to be closed, got %v", winner)
		}
	})

	t.Run("UpdateComplete", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		// Update the winner as processed
		err = st.UpdateComplete(context.Background(), winners[0].BlockNumber)
		if err != nil {
			t.Fatalf("Failed to update winner: %s", err)
		}

		ctx, cancel := context.WithCancel(context.Background())

		winnerChan := st.SubscribeWinners(ctx)
		if err != nil {
			t.Fatalf("Failed to subscribe to winners: %s", err)
		}

		winner := <-winnerChan
		if winner.BlockNumber != winners[1].BlockNumber {
			t.Fatalf("Expected block number %d, got %d", winners[1].BlockNumber, winner.BlockNumber)
		}
		if winner.Winner != winners[1].Winner {
			t.Fatalf("Expected builder address %s, got %s", winners[1].Winner, winner.Winner)
		}

		cancel()

		winner, ok := <-winnerChan
		if ok {
			t.Fatalf("Expected channel to be closed, got %v", winner)
		}
	})

	t.Run("AddSettlement", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		for _, settlement := range settlements {
			err = st.AddSettlement(
				context.Background(),
				settlement.CommitmentIdx,
				settlement.TxHash,
				settlement.BlockNum,
				settlement.Amount,
				settlement.Builder,
				settlement.Type,
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

		settlementChan := st.SubscribeSettlements(ctx)
		if err != nil {
			t.Fatalf("Failed to subscribe to settlements: %s", err)
		}

		for i := 0; i < 6; i++ {
			settlement := <-settlementChan
			if diff := cmp.Diff(settlement, settlements[i]); diff != "" {
				t.Fatalf("Unexpected settlement: (-want +have):\n%s", diff)
			}
		}

		cancel()

		settlement, ok := <-settlementChan
		if ok {
			t.Fatalf("Expected channel to be closed, got %v", settlement)
		}
	})

	t.Run("SettlementInitiated", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		indexes := make([][]byte, 2)
		for i := 0; i < 3; i++ {
			indexes[0] = settlements[2*i].CommitmentIdx
			indexes[1] = settlements[2*i+1].CommitmentIdx

			err = st.SettlementInitiated(
				context.Background(),
				indexes,
				common.HexToHash(fmt.Sprintf("0x%02d", i)),
				uint64(i+1),
			)
			if err != nil {
				t.Fatalf("Failed to initiate settlement: %s", err)
			}
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
		if lastNonce != 3 {
			t.Fatalf("Expected last nonce 3, got %d", lastNonce)
		}

		pendingTxnCount, err := st.PendingTxnCount()
		if err != nil {
			t.Fatalf("Failed to get pending txn count: %s", err)
		}
		if pendingTxnCount != 3 {
			t.Fatalf("Expected pending txn count 3, got %d", pendingTxnCount)
		}
	})

	t.Run("MarkSettlementComplete", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		count, err := st.MarkSettlementComplete(context.Background(), 4)
		if err != nil {
			t.Fatalf("Failed to mark settlement complete: %s", err)
		}
		if count != 6 {
			t.Fatalf("Expected count 6, got %d", count)
		}

		pendingTxnCount, err := st.PendingTxnCount()
		if err != nil {
			t.Fatalf("Failed to get pending txn count: %s", err)
		}
		if pendingTxnCount != 0 {
			t.Fatalf("Expected pending txn count 0, got %d", pendingTxnCount)
		}
	})

	t.Run("stats", func(t *testing.T) {
		st, err := store.NewStore(db)
		if err != nil {
			t.Fatalf("Failed to create store: %s", err)
		}

		stats, err := st.CommitmentStats()
		if err != nil {
			t.Fatalf("Failed to get stats: %s", err)
		}
		if stats.TotalCount != 6 {
			t.Fatalf("Expected total count 6, got %d", stats.TotalCount)
		}
		if stats.RewardCount != 2 {
			t.Fatalf("Expected reward count 2, got %d", stats.RewardCount)
		}
		if stats.SlashCount != 2 {
			t.Fatalf("Expected slash count 2, got %d", stats.SlashCount)
		}
		if stats.SettlementsCompletedCount != 6 {
			t.Fatalf("Expected settlements completed count 6, got %d", stats.SettlementsCompletedCount)
		}

		blockStats, err := st.ProcessedBlocks(2, 0)
		if err != nil {
			t.Fatalf("Failed to get processed blocks: %s", err)
		}
		if len(blockStats) != 1 {
			t.Fatalf("Expected 1 block stats, got %d", len(blockStats))
		}
		block := blockStats[0]
		if block.BlockNumber != winners[0].BlockNumber {
			t.Fatalf("Expected block number %d, got %d", winners[0].BlockNumber, block.BlockNumber)
		}
		if block.Builder != winners[0].Winner {
			t.Fatalf("Expected builder address %s, got %s", winners[0].Winner, block.Builder)
		}
		if block.NoOfCommitments != 3 {
			t.Fatalf("Expected no of commitments 3, got %d", block.NoOfCommitments)
		}
		if block.TotalAmount != 4000000 {
			t.Fatalf("Expected total amount 5000000, got %d", block.TotalAmount)
		}
		if block.NoOfRewards != 1 {
			t.Fatalf("Expected no of rewards 1, got %d", block.NoOfRewards)
		}
		if block.TotalRewards != 2000000 {
			t.Fatalf("Expected total rewards 2000000, got %d", block.TotalRewards)
		}
		if block.NoOfSlashes != 1 {
			t.Fatalf("Expected no of slashes 2, got %d", block.NoOfSlashes)
		}
		if block.TotalSlashes != 1000000 {
			t.Fatalf("Expected total slashes 2000000, got %d", block.TotalSlashes)
		}
		if block.NoOfSettlements != 3 {
			t.Fatalf("Expected no of settlements 3, got %d", block.NoOfSettlements)
		}
	})
}
