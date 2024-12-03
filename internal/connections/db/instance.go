package db

import (
	"database/sql"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	sqldblogger "github.com/simukti/sqldb-logger"
	sqldbzerolog "github.com/simukti/sqldb-logger/logadapter/zerologadapter"

	// Import the MySQL driver
	_ "github.com/go-sql-driver/mysql"

	"golang-boilerplate/internal/config"
)

// NewDB creates and configures a new sqlx.DB instance.
// It connects to the database using the provided configuration,
// enables optional debug logging, and applies connection pool settings.
func NewDB(logger *zerolog.Logger, cfg *config.DB) (*sqlx.DB, error) {
	// Set the name mapper to convert Go struct fields to snake_case.
	sqlx.NameMapper = strcase.ToSnake

	// Create a base database connection.
	sqlDB, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, err
	}

	// Enable debug logging if configured.
	if cfg.Debug == "true" {
		// Create logger instance for database operations.
		dbLogger := logger.With().Str("eventClass", "database").Str("event", "statement").Logger()

		// Set up the log adapter for structured logging.
		logAdapter := sqldbzerolog.New(dbLogger)

		// Set up logger options to customize the log fields and format.
		loggerOptions := []sqldblogger.Option{
			sqldblogger.WithStatementIDFieldname("eventID"),
			sqldblogger.WithTimeFieldname("dbExecTime"),
			sqldblogger.WithConnectionIDFieldname("connID"),
			sqldblogger.WithDurationFieldname("elapsedTime"),
			sqldblogger.WithTimeFormat(sqldblogger.TimeFormatRFC3339Nano),
		}

		// Attach logger options to the SQL driver.
		sqlDB = sqldblogger.OpenDriver(cfg.DSN, sqlDB.Driver(), logAdapter, loggerOptions...)
	}

	// Wrap the base database connection with sqlx for enhanced functionality.
	db := sqlx.NewDb(sqlDB, "mysql")

	// Verify the connection to the database.
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Apply connection pool settings.
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(cfg.MaxConnLifetime) * time.Second)

	return db, nil
}
