package db

import (
	"context"
	"errors"

	"github.com/go-sql-driver/mysql"

	"golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/models"
	utils_mysql "golang-boilerplate/internal/pkg/utils/mysql"
)

type ProductBillerRepository interface {
	Create(ctx context.Context, productBiller *models.ProductBiller) error
}

type productBillerRepository struct {
	db db.DBExecutor
}

func NewProductBillerRepository(db db.DBExecutor) ProductBillerRepository {
	return &productBillerRepository{
		db: db,
	}
}

func (r *productBillerRepository) Create(ctx context.Context, productBiller *models.ProductBiller) error {
	sqlStr := "INSERT INTO product_billers (product_id, biller_id, is_active, created_at, created_by, updated_at, updated_by) "
	sqlStr += "VALUES (:product_id, :biller_id, :is_active, NOW(6), :created_by, NOW(6), :updated_by)"
	_, err := r.db.NamedExecContext(ctx, sqlStr, productBiller)

	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				msg := utils_mysql.ParseDuplicateEntry(err.Error())
				return errors.New(msg)
			}
		}
		return err
	}

	return nil
}
