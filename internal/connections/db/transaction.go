package db

import (
	"context"

	"github.com/jmoiron/sqlx"

	"golang-boilerplate/internal/utils/logger"
)

// WithTransaction handles the lifecycle of a transaction (begin, commit, rollback).
func WithTransaction(ctx context.Context, db DBExecutor, eventClass string, eventName string, fn func(tx *sqlx.Tx) error) (err error) {
	// Check if we're already in a transaction
	if tx, ok := db.(*sqlx.Tx); ok {
		// If already in a transaction, use it directly
		return fn(tx)
	}

	// Start a new transaction
	tx, err := db.(*sqlx.DB).BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				l := logger.FromContext(ctx)
				l.Error(ctx, eventClass, eventName, "tx.Rollback() failed: %v", rollbackErr.Error())
			}
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				l := logger.FromContext(ctx)
				l.Error(ctx, eventClass, eventName, "tx.Rollback() failed: %v", rollbackErr.Error())
			}
		} else {
			err = tx.Commit() // if Commit returns error, update err
		}
	}()

	// Execute the provided function
	err = fn(tx)
	return err
}
