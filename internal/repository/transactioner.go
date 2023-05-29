package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// TransactionManager is a struct that manages transactions
type TransactionManager struct {
	db *sqlx.DB
}

// NewTransactionManager creates a new TransactionManager
func NewTransactionManager(db *sqlx.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

func (t *TransactionManager) Begin(ctx context.Context, o *sql.TxOptions) (*sqlx.Tx, error) {
	tx, err := t.db.BeginTxx(ctx, o)
	if err != nil {
		return nil, errors.Wrap(err, "Error when beginning the transaction")
	}

	return tx, nil
}

func (t *TransactionManager) Commit(ctx context.Context, tx *sqlx.Tx) error {
	return tx.Commit()
}

func (t *TransactionManager) Rollback(ctx context.Context, tx *sqlx.Tx) error {
	err := tx.Rollback()
	if err != nil {
		return errors.Wrap(err, "Error when rolling back the transaction")
	}

	return nil
}
