package repository

import (
	"context"
	"database/sql"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// We can maintain a skipped block list in the smart contract
type PreConfirmationsContract interface {
	GetCommitmentsByBlockNumber(opts *bind.CallOpts, blockNumber *big.Int) ([][32]byte, error)
	GetTxnHashFromCommitment(opts *bind.CallOpts, commitmentIndex [32]byte) (string, error)
}

// CommitmentsStore is an interface that is used to retrieve commitments from the smart contract
// and store them in a local database
type CommitmentsStore interface {
	UpdateCommitmentsForBlockNumber(blockNumber int64) (done chan struct{}, err chan error)
	RetrieveCommitments(blockNumber int64) (map[string]bool, error)

	// Used for restarting the Commitment Store on startup
	LargestStoredBlockNumber() (int64, error)
}

type DBTxnStore struct {
	db            *sql.DB
	preConfClient PreConfirmationsContract
}

func NewDBTxnStore(db *sql.DB, preConfClient PreConfirmationsContract) CommitmentsStore {
	return &DBTxnStore{
		db:            db,
		preConfClient: preConfClient,
	}
}

func (f DBTxnStore) LargestStoredBlockNumber() (int64, error) {
	var largestBlockNumber int64
	err := f.db.QueryRow("SELECT MAX(block_number) FROM committed_transactions").Scan(&largestBlockNumber)
	if err != nil {
		return 0, err
	}
	return largestBlockNumber, nil
}

func (f DBTxnStore) UpdateCommitmentsForBlockNumber(blockNumber int64) (done chan struct{}, errorC chan error) {
	done = make(chan struct{})
	errorC = make(chan error)

	go func(done chan struct{}, errorC chan error) {
		commitments, err := f.preConfClient.GetCommitmentsByBlockNumber(&bind.CallOpts{
			Pending: false,
			From:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
			Context: context.Background(),
		}, big.NewInt(blockNumber))
		if err != nil {
			log.Error().Err(err).Msg("Error getting commitments")
			errorC <- err
			return
		}
		log.Info().Int("block_number", int(blockNumber)).Int("commitments", len(commitments)).Msg("Retrieved commitments")
		for _, commitment := range commitments {
			commitmentTxnHash, err := f.preConfClient.GetTxnHashFromCommitment(&bind.CallOpts{
				Pending: false,
				From:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
				Context: context.Background(),
			}, commitment)
			if err != nil {
				log.Error().Err(err).Msg("Error getting txn hash from commitment")
				errorC <- err
				return
			}

			sqlStatement := `
			INSERT INTO committed_transactions (transaction, block_number)
			VALUES ($1, $2)`
			result, err := f.db.Exec(sqlStatement, commitmentTxnHash, blockNumber)
			if err != nil {
				if err, ok := err.(*pq.Error); ok {
					// Check if the error is a duplicate key violation
					if err.Code.Name() == "unique_violation" {
						log.Info().Msg("Duplicate key violation - ignoring")
						continue
					}
				}
				log.Error().Err(err).Msg("Error inserting txn into DB")
				errorC <- err
				return
			}
			rowsImpacted, err := result.RowsAffected()
			if err != nil {
				log.Error().Err(err).Msg("Error getting rows impacted")
				errorC <- err
				return
			}
			log.Debug().Int("rows_affected", int(rowsImpacted)).Msg("Inserted txn into DB")
		}
		done <- struct{}{}
	}(done, errorC)

	return done, errorC
}

func (f DBTxnStore) RetrieveCommitments(blockNumber int64) (map[string]bool, error) {
	Store := make(map[string]bool)

	rows, err := f.db.Query("SELECT transaction FROM committed_transactions WHERE block_number = $1", blockNumber)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var txnHash string
		err = rows.Scan(&txnHash)
		if err != nil {
			return nil, err
		}
		Store[txnHash] = true
	}

	return Store, nil
}

type InMemoryTxnStore struct {
	db            map[int64][]string // BlockNumber -> TxnHashes
	preConfClient PreConfirmationsContract
}

func NewTransactionCommitmentStore(preConfClient PreConfirmationsContract) CommitmentsStore {
	return &InMemoryTxnStore{
		db:            make(map[int64][]string),
		preConfClient: preConfClient,
	}
}

// Reduant as InMemoryTxnStore is not persisten
func (f InMemoryTxnStore) LargestStoredBlockNumber() (int64, error) {
	var largestBlockNumber int64
	for blockNumber := range f.db {
		if blockNumber > largestBlockNumber {
			largestBlockNumber = blockNumber
		}
	}
	return largestBlockNumber, nil
}

func (f InMemoryTxnStore) UpdateCommitmentsForBlockNumber(blockNumber int64) (done chan struct{}, errorC chan error) {
	done = make(chan struct{})
	errorC = make(chan error)
	if _, ok := f.db[blockNumber]; !ok {
		f.db[blockNumber] = []string{}
	}

	go func(done chan struct{}, errorC chan error) {
		commitments, err := f.preConfClient.GetCommitmentsByBlockNumber(&bind.CallOpts{
			Pending: false,
			From:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
			Context: context.Background(),
		}, big.NewInt(blockNumber))
		if err != nil {
			log.Error().Err(err).Msg("Error getting commitments")
			errorC <- err
			return
		}

		for _, commitment := range commitments {
			commitmentTxnHash, err := f.preConfClient.GetTxnHashFromCommitment(&bind.CallOpts{
				Pending: false,
				From:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
				Context: context.Background(),
			}, commitment)
			if err != nil {
				log.Error().Err(err).Msg("Error getting txn hash from commitment")
				errorC <- err
				return
			}

			f.db[blockNumber] = append(f.db[blockNumber], commitmentTxnHash)
			// Must clear this variable
			done <- struct{}{}
		}
	}(done, errorC)

	return done, errorC
}

// TODO(@ckartik): Adds init Store
// RetrieveCommitments initializes the a goroutine that will fetch all the commitments from the smart contract for the given blockNumber
// and return a channel that will be used to Store out transactions that have already been confirmed
// Need to model the creation of pre-confirmations from a searcher
// NOTE: Need to manage situation where the contracts receive a commitment after the block has been updated to blockNumber
func (f InMemoryTxnStore) RetrieveCommitments(blockNumber int64) (Store map[string]bool, err error) {
	Store = make(map[string]bool)
	for _, txn := range f.db[blockNumber] {
		Store[txn] = true
	}

	return Store, nil
}
