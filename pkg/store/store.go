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
    bid_id BYTEA,
    chainhash BYTEA,
    nonce BIGINT,
    settled BOOLEAN,
    decay_percentage BIGINT,
    settlement_window BIGINT
);`

var encryptedCommitmentsTable = `
CREATE TABLE IF NOT EXISTS encrypted_commitments (
    commitment_index BYTEA PRIMARY KEY,
    committer BYTEA,
    commitment_hash BYTEA,
    commitment_signature BYTEA,
    block_number BIGINT
);`

var winnersTable = `
CREATE TABLE IF NOT EXISTS winners (
	block_number BIGINT PRIMARY KEY,
	builder_address BYTEA,
	settlement_window BIGINT
);`

var bidderTable = `
CREATE TABLE IF NOT EXISTS bidders (
    bidder_address BYTEA,
    settlement_window BIGINT,
    amount NUMERIC(24, 0),
    chainhash BYTEA,
    nonce BIGINT,
    settled BOOLEAN,
    PRIMARY KEY (bidder_address, settlement_window)
);`

var transactionsTable = `
CREATE TABLE IF NOT EXISTS sent_transactions (
    hash BYTEA PRIMARY KEY,
    nonce BIGINT,
    settled BOOLEAN
);`

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) (*Store, error) {
	for _, table := range []string{
		settlementType,
		settlementsTable,
		encryptedCommitmentsTable,
		winnersTable,
		bidderTable,
		transactionsTable,
	} {
		_, err := db.Exec(table)
		if err != nil {
			return nil, err
		}
	}

	return &Store{
		db: db,
	}, nil
}

func (s *Store) RegisterWinner(
	ctx context.Context,
	blockNum int64,
	winner []byte,
	window int64,
) error {
	insertStr := "INSERT INTO winners (block_number, builder_address, settlement_window) VALUES ($1, $2, $3)"

	_, err := s.db.ExecContext(ctx, insertStr, blockNum, winner, window)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetWinner(
	ctx context.Context,
	blockNum int64,
) (updater.Winner, error) {
	winner := updater.Winner{}
	err := s.db.QueryRowContext(
		ctx,
		"SELECT builder_address, settlement_window FROM winners WHERE block_number = $1",
		blockNum,
	).Scan(&winner.Winner, &winner.Window)
	if err != nil {
		return winner, err
	}
	return winner, nil
}

func (s *Store) AddEncryptedCommitment(
	ctx context.Context,
	commitmentIdx []byte,
	committer []byte,
	commitmentHash []byte,
	commitmentSignature []byte,
	blockNum int64,
) error {
	columns := []string{
		"commitment_index",
		"committer",
		"commitment_hash",
		"commitment_signature",
		"block_number",
	}
	values := []interface{}{
		commitmentIdx,
		committer,
		commitmentHash,
		commitmentSignature,
		blockNum,
	}
	placeholder := make([]string, len(values))
	for i := range columns {
		placeholder[i] = fmt.Sprintf("$%d", i+1)
	}

	insertStr := fmt.Sprintf(
		"INSERT INTO encrypted_commitments (%s) VALUES (%s)",
		strings.Join(columns, ", "),
		strings.Join(placeholder, ", "),
	)

	_, err := s.db.ExecContext(ctx, insertStr, values...)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) AddSettlement(
	ctx context.Context,
	commitmentIdx []byte,
	txHash string,
	blockNum int64,
	amount uint64,
	builder []byte,
	bidID []byte,
	settlementType settler.SettlementType,
	decayPercentage int64,
	window int64,
) error {
	columns := []string{
		"commitment_index",
		"transaction",
		"block_number",
		"builder_address",
		"type",
		"amount",
		"bid_id",
		"settled",
		"chainhash",
		"nonce",
		"decay_percentage",
		"settlement_window",
	}
	values := []interface{}{
		commitmentIdx,
		txHash,
		blockNum,
		builder,
		settlementType,
		amount,
		bidID,
		false,
		nil,
		0,
		decayPercentage,
		window,
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

func (s *Store) IsSettled(
	ctx context.Context,
	commitmentIdx []byte,
) (bool, error) {
	var settled bool
	err := s.db.QueryRowContext(
		ctx,
		"SELECT EXISTS(SELECT 1 FROM settlements WHERE commitment_index = $1)",
		commitmentIdx,
	).Scan(&settled)
	if err != nil {
		return false, err
	}

	return settled, nil
}

func (s *Store) SubscribeSettlements(
	ctx context.Context,
	window int64,
) <-chan settler.Settlement {
	resChan := make(chan settler.Settlement)

	go func() {
		defer close(resChan)

		queryStr := `
				SELECT
					commitment_index, transaction, block_number,
					builder_address, amount, bid_id, type, decay_percentage
				FROM settlements
				WHERE settlement_window = $1 AND settled = false AND chainhash IS NULL AND type != 'return'
				ORDER BY block_number ASC`

		results, err := s.db.QueryContext(ctx, queryStr, window)
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
				&s.BidID,
				&s.Type,
				&s.DecayPercentage,
			)
			if err != nil {
				_ = results.Close()
				return
			}

			select {
			case <-ctx.Done():
				_ = results.Close()
				return
			case resChan <- s:
			}
		}

		_ = results.Close()
	}()

	return resChan
}

func (s *Store) BidderRegistered(
	ctx context.Context,
	bidder []byte,
	window int64,
	amount int64,
) error {
	insertStr := "INSERT INTO bidders (bidder_address, settlement_window, amount) VALUES ($1, $2, $3)"

	_, err := s.db.ExecContext(ctx, insertStr, bidder, window, amount)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) SubscribeReturns(ctx context.Context, limit int, window int64) <-chan settler.Return {
	resChan := make(chan settler.Return)

	go func() {
		defer close(resChan)

		queryStr := `
				SELECT bidder_address
				FROM bidders
				WHERE settlement_window = $1`

		results, err := s.db.QueryContext(ctx, queryStr, window)
		if err != nil {
			return
		}

		returns := make([][]byte, 0, limit)

		for results.Next() {
			var r []byte
			err = results.Scan(&r)
			if err != nil {
				_ = results.Close()
				return
			}

			returns = append(returns, r)
			if len(returns) == limit {
				select {
				case <-ctx.Done():
					_ = results.Close()
					return
				case resChan <- settler.Return{
					Window:  window,
					Bidders: returns,
				}:
					returns = returns[:0]
				}
			}
		}

		if len(returns) > 0 {
			select {
			case <-ctx.Done():
				_ = results.Close()
				return
			case resChan <- settler.Return{
				Window:  window,
				Bidders: returns,
			}:
			}
		}

		_ = results.Close()
	}()

	return resChan
}

func (s *Store) SettlementInitiated(
	ctx context.Context,
	commitmentIdx []byte,
	txHash common.Hash,
	nonce uint64,
) error {
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE settlements SET chainhash = $1, nonce = $2 WHERE commitment_index = $3",
		txHash.Bytes(),
		nonce,
		commitmentIdx,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) ReturnInitiated(
	ctx context.Context,
	window int64,
	bidders [][]byte,
	txHash common.Hash,
	nonce uint64,
) error {
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE bidders SET chainhash = $1, nonce = $2 WHERE settlement_window = $3 AND bidder_address = ANY($4)",
		txHash.Bytes(),
		nonce,
		window,
		pq.Array(bidders),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) MarkSettlementComplete(ctx context.Context, nonce uint64) (int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	_, err = tx.ExecContext(
		ctx,
		"UPDATE settlements SET settled = true WHERE settled = false AND nonce < $1 AND chainhash IS NOT NULL",
		nonce,
	)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	_, err = tx.ExecContext(
		ctx,
		"UPDATE bidders SET settled = true WHERE settled = false AND nonce < $1 AND chainhash IS NOT NULL",
		nonce,
	)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	result, err := tx.ExecContext(
		ctx,
		"UPDATE sent_transactions SET settled = true WHERE settled = false AND nonce < $1",
		nonce,
	)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (s *Store) LastNonce() (int64, error) {
	var lastNonce int64
	err := s.db.QueryRow("SELECT MAX(nonce) FROM sent_transactions").Scan(&lastNonce)
	if err != nil {
		return 0, err
	}
	return lastNonce, nil
}

func (s *Store) SentTxn(nonce uint64, txHash common.Hash) error {
	_, err := s.db.Exec(
		"INSERT INTO sent_transactions (hash, nonce, settled) VALUES ($1, $2, false)",
		txHash.Bytes(),
		nonce,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) PendingTxnCount() (int, error) {
	var count int
	err := s.db.QueryRow(
		"SELECT COUNT(*) FROM sent_transactions WHERE hash IS NOT NULL AND settled = false",
	).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
