package db

import (
	"context"
	"fmt"
	"math"

	"github.com/jmoiron/sqlx"
)

type Pagination struct {
	CountRows bool       // Whether to count total rows in the dataset.
	Order     string     // Order by clause for sorting the results.
	Limit     int        // Number of rows per page.
	Page      int        // Current page number.
	TotalRow  int        // Total number of rows (set if CountRows is true).
	TotalPage int        // Total number of pages (calculated from TotalRow and Limit).
	Rows      *sqlx.Rows // Result set for the paginated query.
}

// Returns the SQL ORDER BY clause if an order is specified.
func (p *Pagination) GetOrder() string {
	if p.Order != "" {
		return "ORDER BY " + p.Order
	}
	return ""
}

// Ensures a valid page is returned (default is 1 if not set).
func (p *Pagination) GetPage() int {
	if p.Page <= 0 { // Guard against zero or negative values.
		p.Page = 1
	}
	return p.Page
}

// Ensures a valid limit is returned (default is 10 if not set).
func (p *Pagination) GetLimit() int {
	if p.Limit <= 0 { // Guard against zero or negative values.
		p.Limit = 10
	}
	return p.Limit
}

// Calculates the offset based on the current page and limit.
func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

// Executes a paginated query with optional row count.
func Paginate(ctx context.Context, pagination *Pagination, db DBExecutor, baseSQL string, args []interface{}) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	// Construct the paginated SQL query.
	sqlPaginated := fmt.Sprintf(
		"%s %s LIMIT %d OFFSET %d",
		baseSQL,
		pagination.GetOrder(),
		pagination.GetLimit(),
		pagination.GetOffset(),
	)

	// Count total rows if required.
	if pagination.CountRows {
		sqlCount := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count", baseSQL)
		if err := db.QueryRowxContext(ctx, sqlCount, args...).Scan(&pagination.TotalRow); err != nil {
			return fmt.Errorf("failed to count rows: %w", err)
		}
		pagination.TotalPage = int(math.Ceil(float64(pagination.TotalRow) / float64(pagination.GetLimit())))
	}

	// Execute the paginated query and retrieve the rows.
	rows, err := db.QueryxContext(ctx, sqlPaginated, args...)
	if err != nil {
		return fmt.Errorf("failed to execute paginated query: %w", err)
	}
	pagination.Rows = rows
	return nil
}
