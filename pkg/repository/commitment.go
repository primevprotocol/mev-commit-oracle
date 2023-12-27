package repository

import (
	"context"
	"database/sql"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/primevprotocol/mev-oracle/pkg/preconf"
	"github.com/rs/zerolog/log"
)

// We can maintain a skipped block list in the smart contract
type PreConfirmationsContract interface {
	GetCommitmentsByBlockNumber(opts *bind.CallOpts, blockNumber *big.Int) ([][32]byte, error)
	GetTxnHashFromCommitment(opts *bind.CallOpts, commitmentIndex [32]byte) (string, error)
	GetCommitment(opts *bind.CallOpts, commitmentIndex [32]byte) (preconf.PreConfCommitmentStorePreConfCommitment, error)
}

// CommitmentsStore is an interface that is used to retrieve commitments from the smart contract
// and store them in a local database
type CommitmentsStore interface {
	UpdateCommitmentsForBlockNumber(blockNumber int64) (done chan struct{}, err chan error)
	RetrieveCommitments(blockNumber int64) ([]Commitment, error)

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
		commitmentIndexes, err := f.preConfClient.GetCommitmentsByBlockNumber(&bind.CallOpts{
			Pending: false,
			From:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
			Context: context.Background(),
		}, big.NewInt(blockNumber))
		if err != nil {
			log.Error().Err(err).Msg("Error getting commitments")
			errorC <- err
			return
		}
		log.Info().Int("block_number", int(blockNumber)).Int("commitments", len(commitmentIndexes)).Msg("Retrieved commitment indexes")
		for _, commitmentIndex := range commitmentIndexes {
			commitment, err := f.preConfClient.GetCommitment(&bind.CallOpts{
				Pending: false,
				From:    common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"),
				Context: context.Background(),
			}, commitmentIndex)
			if err != nil {
				log.Error().Err(err).Msg("Error getting txn hash from commitment")
				errorC <- err
				return
			}
			//sqlStatement := `
			//INSERT INTO committed_transactions (transaction, block_number)
			//VALUES ($1, $2)`

			insertSqlStatement := `
			INSERT INTO committed_transactions (commitment_index, transaction, block_number, builder_address)
			VALUES ($1, $2, $3, $4)`
			result, err := f.db.Exec(insertSqlStatement, commitmentIndex, commitment.TxnHash, commitment.BlockNumber, commitment.Commiter.Bytes())
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

type Commitment struct {
	CommitmentIndex [32]byte
	TxnHash         string
	BlockNum        int64
	BuilderAddress  [32]byte // Assuming builder_address is stored in binary format
}

func (f DBTxnStore) RetrieveCommitments(blockNumber int64) ([]Commitment, error) {
	commitments := make([]Commitment, 0)

	rows, err := f.db.Query("SELECT commitment_index, transaction, block_number, builder_address FROM committed_transactions WHERE block_number = $1", blockNumber)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			commitmentIndex [32]byte
			txnHash         string
			blockNum        int64
			builderAddress  [32]byte // Assuming builder_address is stored in binary format
		)

		// Scan all the columns returned by the query
		err = rows.Scan(&commitmentIndex, &txnHash, &blockNum, &builderAddress)
		if err != nil {
			return nil, err
		}

		commitments = append(commitments, Commitment{
			CommitmentIndex: commitmentIndex,
			TxnHash:         txnHash,
			BlockNum:        blockNum,
			BuilderAddress:  builderAddress,
		})
	}

	return commitments, nil
}
