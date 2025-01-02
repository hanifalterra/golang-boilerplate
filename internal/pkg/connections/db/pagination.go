package database

import (
	"context"
	"fmt"
	"math"
)

// Pagination holds parameters and metadata for paginated queries.
type Pagination struct {
	Order      string `json:"order"`       // Order by clause for sorting the results.
	Limit      int    `json:"limit"`       // Number of rows per page.
	Page       int    `json:"page"`        // Current page number.
	TotalRows  int    `json:"total_rows"`  // Total number of rows (if CountRows is true).
	TotalPages int    `json:"total_pages"` // Total number of pages (calculated from TotalRow and Limit).
}

// GetOrder constructs the SQL ORDER BY clause.
func (p *Pagination) GetOrder() string {
	if p.Order != "" {
		return "ORDER BY " + p.Order
	}
	return ""
}

// GetPage ensures a valid page is returned.
func (p *Pagination) GetPage() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return p.Page
}

// GetLimit ensures a valid limit is returned.
func (p *Pagination) GetLimit() int {
	if p.Limit <= 0 {
		p.Limit = 10
	}
	return p.Limit
}

// GetOffset calculates the offset for the SQL query.
func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

// Paginate executes a paginated query with optional row counting and populates the results slice.
func Paginate[T any](
	ctx context.Context,
	db DBExecutor,
	baseSQL string,
	args []interface{},
	pagination *Pagination,
	results *[]T,
) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	// Count total rows.
	sqlCount := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count", baseSQL)
	if err := db.QueryRowxContext(ctx, sqlCount, args...).Scan(&pagination.TotalRows); err != nil {
		return fmt.Errorf("failed to count rows: %w", err)
	}
	pagination.TotalPages = int(math.Ceil(float64(pagination.TotalRows) / float64(pagination.GetLimit())))

	// Construct the paginated SQL query.
	sqlPaginated := fmt.Sprintf(
		"%s %s LIMIT %d OFFSET %d",
		baseSQL,
		pagination.GetOrder(),
		pagination.GetLimit(),
		pagination.GetOffset(),
	)

	// Execute the paginated query.
	rows, err := db.QueryxContext(ctx, sqlPaginated, args...)
	if err != nil {
		return fmt.Errorf("failed to execute paginated query: %w", err)
	}
	defer rows.Close()

	// Populate the results.
	for rows.Next() {
		var result T
		if err := rows.StructScan(&result); err != nil {
			return fmt.Errorf("failed to scan row into struct: %w", err)
		}
		*results = append(*results, result)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("row iteration error: %w", err)
	}

	return nil
}
