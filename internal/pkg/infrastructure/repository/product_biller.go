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
	Update(ctx context.Context, productID, billerID uint, productBiller *models.ProductBiller) error
	Delete(ctx context.Context, productID, billerID uint) error
	GetOne(ctx context.Context, productID, billerID uint) (*models.ProductBiller, error)
	GetMany(ctx context.Context, filter map[string]interface{}) ([]*models.ProductBiller, error)
	GetManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.ProductBiller, *db.Pagination, error)
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

func (r *productBillerRepository) Create(ctx context.Context, productBiller *models.ProductBiller) error {
	const query = `
		INSERT INTO product_billers 
		(product_id, biller_id, is_active, created_at, created_by, updated_at, updated_by)
		VALUES (:product_id, :biller_id, :is_active, NOW(6), :created_by, NOW(6), :updated_by)
	`

	_, err := r.db.NamedExecContext(ctx, query, productBiller)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			duplicateField := utils_db.ParseDuplicateEntry(mysqlErr.Message)
			return fmt.Errorf("duplicate entry detected: %s", duplicateField)
		}
		return fmt.Errorf("failed to create product biller: %w", err)
	}

	return nil
}

func (r *productBillerRepository) Update(ctx context.Context, productID, billerID uint, productBiller *models.ProductBiller) error {
	const query = `
		UPDATE product_billers
		SET is_active = :is_active, updated_at = NOW(6)
		WHERE product_id = :product_id AND biller_id = :biller_id AND deleted_at IS NULL
	`

	params := map[string]interface{}{
		"product_id": productID,
		"biller_id":  billerID,
		"is_active":  productBiller.IsActive,
	}

	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update product biller: %w", err)
	}

	return nil
}

func (r *productBillerRepository) Delete(ctx context.Context, productID, billerID uint) error {
	const query = `
		UPDATE product_billers
		SET deleted_at = NOW(6)
		WHERE product_id = :product_id AND biller_id = :biller_id AND deleted_at IS NULL
	`

	params := map[string]interface{}{
		"product_id": productID,
		"biller_id":  billerID,
	}

	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to delete product biller: %w", err)
	}

	return nil
}

func (r *productBillerRepository) getBaseQuery(filter map[string]interface{}) (string, []interface{}) {
	const baseQuery = `
		SELECT product_id, biller_id, is_active, created_at, created_by, updated_at, updated_by
		FROM product_billers
		WHERE deleted_at IS NULL
	`

	return utils_db.ApplyFilters(baseQuery, filter)
}

func (r *productBillerRepository) GetOne(ctx context.Context, productID, billerID uint) (*models.ProductBiller, error) {
	filter := map[string]interface{}{
		"product_id": productID,
		"biller_id":  billerID,
	}

	query, args := r.getBaseQuery(filter)

	var productBiller models.ProductBiller
	if err := r.db.GetContext(ctx, &productBiller, query, args...); err != nil {
		return nil, fmt.Errorf("failed to fetch product biller: %w", err)
	}

	return &productBiller, nil
}

func (r *productBillerRepository) GetMany(ctx context.Context, filter map[string]interface{}) ([]*models.ProductBiller, error) {
	query, args := r.getBaseQuery(filter)

	var productBillers []*models.ProductBiller
	if err := r.db.SelectContext(ctx, &productBillers, query, args...); err != nil {
		return nil, fmt.Errorf("failed to fetch product billers: %w", err)
	}

	return productBillers, nil
}

func (r *productBillerRepository) GetManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.ProductBiller, *db.Pagination, error) {
	query, args := r.getBaseQuery(filter)

	pagination := &db.Pagination{Order: "id ASC", Page: page, Limit: limit}
	var productBillers []*models.ProductBiller
	if err := db.Paginate(ctx, r.db, query, args, pagination, &productBillers); err != nil {
		return nil, nil, fmt.Errorf("failed to fetch product billers with pagination: %w", err)
	}

	return productBillers, pagination, nil
}
