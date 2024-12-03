package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// DBExecutor defines an abstraction for database operations compatible with *sqlx.DB and *sqlx.Tx.
// It ensures all methods support context, promoting better control over timeouts and cancellations.
type DBExecutor interface {
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	DriverName() string
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
	Rebind(query string) string
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}
