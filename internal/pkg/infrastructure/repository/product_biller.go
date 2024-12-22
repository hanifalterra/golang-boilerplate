package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"

	db "golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/models"
	utils_db "golang-boilerplate/internal/pkg/utils/db"
)

// ProductBillerRepository defines the interface for managing ProductBiller entities.
type ProductBillerRepository interface {
	Create(ctx context.Context, productBiller *models.ProductBiller) error
	Update(ctx context.Context, productID, billerID int, productBiller *models.ProductBiller) error
	Delete(ctx context.Context, productID, billerID int) error
	GetOne(ctx context.Context, productID, billerID int) (*models.ProductBiller, error)
	GetMany(ctx context.Context, filter map[string]interface{}) ([]*models.ProductBiller, error)
	GetManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.ProductBiller, int, error)
}

// productBillerRepository implements ProductBillerRepository.
type productBillerRepository struct {
	db db.DBExecutor
}

// NewProductBillerRepository creates a new instance of ProductBillerRepository.
func NewProductBillerRepository(db db.DBExecutor) ProductBillerRepository {
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
				duplicateField := utils_db.ParseDuplicateEntry(mysqlErr.Message)
				return fmt.Errorf("duplicate entry detected: %s", duplicateField)
			}
		}
		// Return a wrapped error with context
		return fmt.Errorf("failed to create product biller: %w", err)
	}

	return nil
}

// UpdateIsActive updates the `is_active` field for a given product and biller.
func (r *productBillerRepository) Update(ctx context.Context, productID, billerID int, productBiller *models.ProductBiller) error {
	const query = `
		UPDATE product_billers
		SET is_active = :is_active, updated_at = NOW(6)
		WHERE product_id = :product_id AND biller_id = :biller_id AND is_deleted = false
	`

	params := map[string]interface{}{
		"product_id": productID,
		"biller_id":  billerID,
		"is_active":  productBiller.IsActive,
	}

	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update is_active for product biller: %w", err)
	}

	return nil
}

func (r *productBillerRepository) Delete(ctx context.Context, productID, billerID int) error {
	const query = `
		UPDATE product_billers
		SET is_deleted = true, updated_at = NOW(6)
		WHERE product_id = :product_id AND biller_id = :biller_id AND is_deleted = false
	`

	params := map[string]interface{}{
		"product_id": productID,
		"biller_id":  billerID,
	}

	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to soft delete product biller: %w", err)
	}

	return nil
}

// GetOne retrieves a single ProductBiller based on product and biller IDs.
func (r *productBillerRepository) GetOne(ctx context.Context, productID, billerID int) (*models.ProductBiller, error) {
	const query = `
		SELECT product_id, biller_id, is_active, created_at, created_by, updated_at, updated_by
		FROM product_billers
		WHERE product_id = :product_id AND biller_id = :biller_id AND is_deleted = false
	`

	params := map[string]interface{}{
		"product_id": productID,
		"biller_id":  billerID,
	}

	var productBiller models.ProductBiller
	if err := r.db.GetContext(ctx, &productBiller, query, params); err != nil {
		return nil, fmt.Errorf("failed to fetch product biller: %w", err)
	}

	return &productBiller, nil
}

// getBaseQuery builds the base query for fetching ProductBiller records with filters.
func (r *productBillerRepository) getBaseQuery(filter map[string]interface{}) (string, []interface{}) {
	const baseQuery = `
		SELECT product_id, biller_id, is_active, created_at, created_by, updated_at, updated_by
		FROM product_billers
		WHERE is_deleted = false
	`

	// Use utils_mysql.ApplyFilters to add dynamic filters
	query, args := utils_db.ApplyFilters(baseQuery, filter)
	return query, args
}

// GetMany retrieves multiple ProductBillers based on a set of filters.
func (r *productBillerRepository) GetMany(ctx context.Context, filter map[string]interface{}) ([]*models.ProductBiller, error) {
	query, args := r.getBaseQuery(filter)

	var productBillers []*models.ProductBiller
	if err := r.db.SelectContext(ctx, &productBillers, query, args...); err != nil {
		return nil, fmt.Errorf("failed to fetch product billers: %w", err)
	}

	return productBillers, nil
}

// GetManyWithPagination retrieves multiple ProductBillers with filters and pagination.
func (r *productBillerRepository) GetManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.ProductBiller, int, error) {
	query, args := r.getBaseQuery(filter)

	// Add pagination to the query
	offset := (page - 1) * limit
	query = fmt.Sprintf("%s LIMIT %d OFFSET %d", query, limit, offset)

	// Fetch total count for pagination
	countQuery := `
		SELECT COUNT(*) AS total
		FROM product_billers
		WHERE is_deleted = false
	`
	countQuery, countArgs := utils_db.ApplyFilters(countQuery, filter)

	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, countArgs...); err != nil {
		return nil, 0, fmt.Errorf("failed to count product billers: %w", err)
	}

	// Fetch paginated data
	var productBillers []*models.ProductBiller
	if err := r.db.SelectContext(ctx, &productBillers, query, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to fetch product billers with pagination: %w", err)
	}

	return productBillers, total, nil
}
