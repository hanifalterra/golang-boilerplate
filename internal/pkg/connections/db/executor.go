package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// DBExecutor abstracts database operations compatible with *sqlx.DB and *sqlx.Tx.
// This interface ensures all methods support context for better timeout and cancellation control.
// It includes common methods for querying, executing, and preparing statements.
type DBExecutor interface {
	// BindNamed binds named parameters in a query, returning the bound query and arguments.
	BindNamed(query string, arg interface{}) (string, []interface{}, error)

	// DriverName returns the name of the driver used for database operations.
	DriverName() string

	// ExecContext executes a query that doesn't return rows (e.g., INSERT, UPDATE).
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)

	// GetContext retrieves a single result and populates the dest with the result.
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error

	// NamedExecContext executes a query with named parameters.
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)

	// PrepareContext prepares a statement for execution. Must be closed when no longer needed.
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)

	// PrepareNamedContext prepares a named statement for execution. Must be closed when no longer needed.
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)

	// PreparexContext prepares an extended SQL statement for execution. Must be closed when no longer needed.
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)

	// QueryContext executes a query and returns rows.
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

	// QueryRowContext executes a query that returns a single row.
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row

	// QueryRowxContext executes a query that returns a single row with extended capabilities.
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row

	// QueryxContext executes a query and returns extended rows.
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)

	// Rebind rebinds the query for the driver's parameter format.
	Rebind(query string) string

	// SelectContext executes a query and populates dest with the result set.
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// Interface checker
var (
	_ DBExecutor = (*sqlx.DB)(nil)
	_ DBExecutor = (*sqlx.Tx)(nil)
)
