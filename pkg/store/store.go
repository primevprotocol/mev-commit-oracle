package store

import (
	"context"
	"database/sql"

	"github.com/ethereum/go-ethereum/common"
	"github.com/primevprotocol/mev-oracle/pkg/settler"
	"github.com/primevprotocol/mev-oracle/pkg/updater"
)

var settlementsTable = `
CREATE TABLE IF NOT EXISTS settlements (
    commitment_index BYTEA PRIMARY KEY,
    transaction VARCHAR(255),
    block_number BIGINT,
    builder_address BYTEA,
    is_slash BOOLEAN,
    nonce BIGINT,
    chainhash BYTEA,
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
	_, err := db.Exec(settlementsTable)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(winnersTable)
	if err != nil {
		return nil, err
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
				resChan <- bWinner
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
	builder string,
	isSlash bool,
) error {
	insertStr := `
		INSERT INTO settlements (commitment_index, transaction, block_number, builder_address, is_slash, nonce, chainhash, settled)
		VALUES ($1, $2, $3, $4, $5, 0, NULL, false)`

	_, err := s.db.ExecContext(ctx, insertStr, commitmentIdx, txHash, blockNum, builder, isSlash)
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
				SELECT commitment_index, transaction, block_number, builder_address, is_slash
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
					&s.IsSlash,
				)
				if err != nil {
					continue RETRY
				}

				resChan <- s
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
	commitmentIdx []byte,
	txHash common.Hash,
	nonce uint64,
) error {
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE settlements SET transaction = $1, nonce = $2 WHERE commitment_index = $3",
		txHash.Bytes(),
		nonce,
		commitmentIdx,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) MarkSettlementComplete(ctx context.Context, nonce uint64) error {
	_, err := s.db.ExecContext(
		ctx,
		"UPDATE settlements SET settled = true WHERE nonce < $1 AND chainhash IS NOT NULL",
		nonce,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) LastNonce() (int64, error) {
	var lastNonce int64
	err := s.db.QueryRow("SELECT MAX(nonce) FROM settlements").Scan(&lastNonce)
	if err != nil {
		return 0, err
	}
	return lastNonce, nil
}
