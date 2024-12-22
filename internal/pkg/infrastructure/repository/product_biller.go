package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"

	"golang-boilerplate/internal/pkg/connections/database"
	"golang-boilerplate/internal/pkg/models"
	utils_mysql "golang-boilerplate/internal/pkg/utils/mysql"
)

// ProductBillerRepository defines the interface for managing ProductBiller entities.
type ProductBillerRepository interface {
	Create(ctx context.Context, productBiller *models.ProductBiller) error
}

// productBillerRepository implements ProductBillerRepository.
type productBillerRepository struct {
	db database.DBExecutor
}

// NewProductBillerRepository creates a new instance of ProductBillerRepository.
func NewProductBillerRepository(db database.DBExecutor) ProductBillerRepository {
	return &productBillerRepository{
		db: db,
	}
}

// Create inserts a new ProductBiller record into the database.
func (r *productBillerRepository) Create(ctx context.Context, productBiller *models.ProductBiller) error {
	const query = `
		INSERT INTO product_billers 
		(product_id, biller_id, is_active, created_at, created_by, updated_at, updated_by)
		VALUES (:product_id, :biller_id, :is_active, NOW(6), :created_by, NOW(6), :updated_by)
	`

	_, err := r.db.NamedExecContext(ctx, query, productBiller)
	if err != nil {
		// Check if the error is a MySQL-specific error
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				// Parse and return a formatted duplicate entry error
				duplicateField := utils_mysql.ParseDuplicateEntry(mysqlErr.Message)
				return fmt.Errorf("duplicate entry detected: %s", duplicateField)
			}
		}
		// Return a wrapped error with context
		return fmt.Errorf("failed to create product biller: %w", err)
	}

	return nil
}
