package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/primevprotocol/mev-oracle/pkg/settler"
	"github.com/primevprotocol/mev-oracle/pkg/updater"
)

var settlementType = `
DO $$ BEGIN
    CREATE TYPE settlement_type AS ENUM ('reward', 'slash', 'return');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;`

var settlementsTable = `
CREATE TABLE IF NOT EXISTS settlements (
    commitment_index BYTEA PRIMARY KEY,
    transaction TEXT,
    block_number BIGINT,
    builder_address BYTEA,
    type settlement_type,
    amount NUMERIC(24, 0),
    chainhash BYTEA,
    nonce BIGINT,
    settled BOOLEAN
);`

var winnersTable = `
CREATE TABLE IF NOT EXISTS winners (
    block_number BIGINT PRIMARY KEY,
    builder_address BYTEA,
    processed BOOLEAN
);`

type Store struct {
	db       *sql.DB
	winnerT  chan struct{}
	settlerT chan struct{}
}

func NewStore(db *sql.DB) (*Store, error) {
	for _, table := range []string{settlementType, settlementsTable, winnersTable} {
		_, err := db.Exec(table)
		if err != nil {
			return nil, err
		}
	}

	return &Store{
		db:       db,
		winnerT:  make(chan struct{}),
		settlerT: make(chan struct{}),
	}, nil
}

func (s *Store) triggerWinner() {
	select {
	case s.winnerT <- struct{}{}:
	default:
	}
}

func (s *Store) triggerSettler() {
	select {
	case s.settlerT <- struct{}{}:
	default:
	}
}

func (s *Store) RegisterWinner(ctx context.Context, blockNum int64, winner string) error {
	insertStr := "INSERT INTO winners (block_number, builder_address, processed) VALUES ($1, $2, $3)"

	_, err := s.db.ExecContext(ctx, insertStr, blockNum, winner, false)
	if err != nil {
		return err
	}
	s.triggerWinner()
	return nil
}

func (s *Store) SubscribeWinners(ctx context.Context) <-chan updater.BlockWinner {
	resChan := make(chan updater.BlockWinner)
	go func() {
		defer close(resChan)

	RETRY:
		for {
			results, err := s.db.QueryContext(
				ctx,
				"SELECT block_number, builder_address FROM winners WHERE processed = false",
			)
			if err != nil {
				return
			}
			for results.Next() {
				var bWinner updater.BlockWinner
				err = results.Scan(&bWinner.BlockNumber, &bWinner.Winner)
				if err != nil {
					continue RETRY
				}
				select {
				case <-ctx.Done():
					return
				case resChan <- bWinner:
				}
			}

			select {
			case <-ctx.Done():
				return
			case <-s.winnerT:
			}
		}
	}()

	return resChan
}

func (s *Store) UpdateComplete(ctx context.Context, blockNum int64) error {
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE winners SET processed = true WHERE block_number = $1",
		blockNum,
	)
	if err != nil {
		return err
	}
	s.triggerSettler()
	return nil
}

func (s *Store) AddSettlement(
	ctx context.Context,
	commitmentIdx []byte,
	txHash string,
	blockNum int64,
	amount uint64,
	builder string,
	settlementType settler.SettlementType,
) error {
	columns := []string{
		"commitment_index",
		"transaction",
		"block_number",
		"builder_address",
		"type",
		"amount",
		"settled",
		"chainhash",
		"nonce",
	}
	values := []interface{}{
		commitmentIdx,
		txHash,
		blockNum,
		builder,
		settlementType,
		amount,
		false,
		nil,
		0,
	}
	placeholder := make([]string, len(values))
	for i := range columns {
		placeholder[i] = fmt.Sprintf("$%d", i+1)
	}

	insertStr := fmt.Sprintf(
		"INSERT INTO settlements (%s) VALUES (%s)",
		strings.Join(columns, ", "),
		strings.Join(placeholder, ", "),
	)

	_, err := s.db.ExecContext(ctx, insertStr, values...)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) SubscribeSettlements(ctx context.Context) <-chan settler.Settlement {
	resChan := make(chan settler.Settlement)

	go func() {
		defer close(resChan)

	RETRY:
		for {
			queryStr := `
				SELECT commitment_index, transaction, block_number, builder_address, amount, type
				FROM settlements
				WHERE settled = false AND chainhash IS NULL`

			results, err := s.db.QueryContext(ctx, queryStr)
			if err != nil {
				return
			}

			for results.Next() {
				var s settler.Settlement

				err = results.Scan(
					&s.CommitmentIdx,
					&s.TxHash,
					&s.BlockNum,
					&s.Builder,
					&s.Amount,
					&s.Type,
				)
				if err != nil {
					continue RETRY
				}

				select {
				case <-ctx.Done():
					return
				case resChan <- s:
				}
			}

			select {
			case <-ctx.Done():
				return
			case <-s.settlerT:
			}
		}
	}()

	return resChan
}

func (s *Store) SettlementInitiated(
	ctx context.Context,
	commitmentIndexes [][]byte,
	txHash common.Hash,
	nonce uint64,
) error {
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE settlements SET chainhash = $1, nonce = $2 WHERE commitment_index = ANY($3::BYTEA[])",
		txHash.Bytes(),
		nonce,
		pq.Array(commitmentIndexes),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) MarkSettlementComplete(ctx context.Context, nonce uint64) (int, error) {
	result, err := s.db.ExecContext(
		ctx,
		"UPDATE settlements SET settled = true WHERE settled = false AND nonce < $1 AND chainhash IS NOT NULL",
		nonce,
	)
	if err != nil {
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (s *Store) LastNonce() (int64, error) {
	var lastNonce int64
	err := s.db.QueryRow("SELECT MAX(nonce) FROM settlements").Scan(&lastNonce)
	if err != nil {
		return 0, err
	}
	return lastNonce, nil
}

func (s *Store) PendingTxnCount() (int, error) {
	var count int
	err := s.db.QueryRow(
		"SELECT COUNT(DISTINCT chainhash) FROM settlements WHERE chainhash IS NOT NULL AND settled = false",
	).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

type BlockInfo struct {
	BlockNumber     int64
	Builder         string
	NoOfCommitments int
	TotalAmount     int
	NoOfRewards     int
	TotalRewards    int
	NoOfSlashes     int
	TotalSlashes    int
	NoOfSettlements int
}

func (s *Store) ProcessedBlocks(limit, offset int) ([]BlockInfo, error) {
	var blocks []BlockInfo
	rows, err := s.db.Query(`
		SELECT
			winners.block_number,
			winners.builder_address,
			COUNT(settlements.commitment_index) AS commitment_count,
			SUM(settlements.amount) AS total_amount,
			COUNT(settlements.type = 'reward' OR NULL) AS reward_count,
			SUM(settlements.amount) FILTER (WHERE settlements.type = 'reward') AS total_rewards,
			COUNT(settlements.type = 'slash' OR NULL) AS slash_count,
			SUM(settlements.amount) FILTER (WHERE settlements.type = 'slash') AS total_slashes,
			COUNT(settlements.settled) AS settled_count
		FROM
			winners
		LEFT JOIN
			settlements ON settlements.block_number = winners.block_number
		WHERE
			winners.processed = true
		GROUP BY
			winners.block_number, winners.builder_address
		ORDER BY
			winners.block_number DESC
		LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var b BlockInfo
		err := rows.Scan(
			&b.BlockNumber,
			&b.Builder,
			&b.NoOfCommitments,
			&b.TotalAmount,
			&b.NoOfRewards,
			&b.TotalRewards,
			&b.NoOfSlashes,
			&b.TotalSlashes,
			&b.NoOfSettlements,
		)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, b)
	}
	return blocks, nil
}

type CommitmentStats struct {
	TotalCount                int
	RewardCount               int
	SlashCount                int
	SettlementsCompletedCount int
}

func (s *Store) CommitmentStats() (CommitmentStats, error) {
	var stats CommitmentStats
	err := s.db.QueryRow(`
		SELECT
			COUNT(*),
			COUNT(type = 'reward' OR NULL),
			COUNT(type = 'slash' OR NULL),
			COUNT(settled)
		FROM
			settlements
	`).Scan(
		&stats.TotalCount,
		&stats.RewardCount,
		&stats.SlashCount,
		&stats.SettlementsCompletedCount,
	)
	if err != nil {
		return stats, err
	}
	return stats, nil
}
