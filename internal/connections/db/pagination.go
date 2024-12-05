package db

import (
	"context"
	"fmt"
	"math"

	"github.com/jmoiron/sqlx"
)

type Pagination struct {
	CountRows bool
	Order     string
	Limit     int
	Page      int
	TotalRow  int
	TotalPage int
	Rows      *sqlx.Rows
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetOrder() string {
	if p.Order != "" {
		return "ORDER BY " + p.Order
	}
	return ""
}

func Paginate(ctx context.Context, pagination *Pagination, db DBExecutor, sql string, args []interface{}) error {
	if db == nil {
		return fmt.Errorf("db is nil")
	}

	sqlPaginated := fmt.Sprintf("%s %s LIMIT %d OFFSET %d", sql, pagination.GetOrder(), pagination.GetLimit(), pagination.GetOffset())

	// Count total rows if required
	if pagination.CountRows {
		sqlCount := fmt.Sprintf("SELECT COUNT(*) FROM (%s) as count", sql)
		err := db.QueryRowxContext(ctx, sqlCount, args...).Scan(&pagination.TotalRow)
		if err != nil {
			return err
		}
		pagination.TotalPage = int(math.Ceil(float64(pagination.TotalRow) / float64(pagination.Limit)))
	}

	// Execute the paginated query
	rows, err := db.QueryxContext(ctx, sqlPaginated, args...)
	if err != nil {
		return err
	}
	pagination.Rows = rows
	return nil
}
