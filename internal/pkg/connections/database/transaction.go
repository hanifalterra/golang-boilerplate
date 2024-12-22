package database

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"golang-boilerplate/internal/pkg/utils/logger"
)

// WithTransaction manages the lifecycle of a database transaction, including commit, rollback, and error propagation.
// If already within a transaction, the existing one is reused.
func WithTransaction(ctx context.Context, db DBExecutor, eventClass, eventName string, fn func(tx *sqlx.Tx) error) (err error) {
	// Reuse the existing transaction if already in one.
	if tx, ok := db.(*sqlx.Tx); ok {
		return fn(tx)
	}

	// Begin a new transaction.
	tx, err := db.(*sqlx.DB).BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "failed to begin transaction")
	}

	// Ensure proper rollback or commit on function exit.
	defer func() {
		// Handle panic scenarios.
		if p := recover(); p != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logger.FromContext(ctx).Error(ctx, eventClass, eventName, "transaction rollback failed during panic: %v", rollbackErr)
			}
			panic(p) // Re-throw panic after rollback.
		}

		// Handle normal flow with error propagation.
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				logger.FromContext(ctx).Error(ctx, eventClass, eventName, "transaction rollback failed: %v", rollbackErr)
				err = errors.Wrapf(err, "transaction failed and rollback also failed: %v", rollbackErr)
			}
		} else {
			// Commit the transaction.
			if commitErr := tx.Commit(); commitErr != nil {
				err = errors.Wrap(commitErr, "failed to commit transaction")
			}
		}
	}()

	// Execute the transactional function.
	err = fn(tx)
	return err
}
