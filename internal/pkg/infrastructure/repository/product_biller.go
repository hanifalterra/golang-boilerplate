package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"

	db "golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/models"
	utils_db "golang-boilerplate/internal/pkg/utils/db"
)

// ProductBillerRepository defines the interface for managing ProductBiller entities.
type ProductBillerRepository interface {
	Create(ctx context.Context, productBiller *models.ProductBiller) error
	Update(ctx context.Context, id uint, productBiller *models.ProductBiller) error
	Delete(ctx context.Context, id uint) error
	DeleteByProductID(ctx context.Context, productID uint) error
	DeleteByBillerID(ctx context.Context, billerID uint) error
	GetOne(ctx context.Context, id uint) (*models.ProductBiller, error)
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

func (r *productBillerRepository) Update(ctx context.Context, id uint, productBiller *models.ProductBiller) error {
	const query = `
		UPDATE product_billers
		SET is_active = :is_active, updated_at = NOW(6)
		WHERE id = :id AND deleted_at IS NULL
	`

	params := map[string]interface{}{
		"id":        id,
		"is_active": productBiller.IsActive,
	}

	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update product biller: %w", err)
	}

	return nil
}

func (r *productBillerRepository) Delete(ctx context.Context, id uint) error {
	const query = `
		UPDATE product_billers
		SET deleted_at = NOW(6)
		WHERE id = :id AND deleted_at IS NULL
	`

	params := map[string]interface{}{
		"id": id,
	}

	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to delete product biller: %w", err)
	}

	return nil
}

func (r *productBillerRepository) DeleteByProductID(ctx context.Context, productID uint) error {
	const query = `
		UPDATE product_billers
		SET deleted_at = NOW(6)
		WHERE product_id = :product_id AND deleted_at IS NULL
	`

	params := map[string]interface{}{
		"product_id": productID,
	}

	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to delete product biller: %w", err)
	}

	return nil
}

func (r *productBillerRepository) DeleteByBillerID(ctx context.Context, billerID uint) error {
	const query = `
		UPDATE product_billers
		SET deleted_at = NOW(6)
		WHERE biller_id = :biller_id AND deleted_at IS NULL
	`

	params := map[string]interface{}{
		"biller_id": billerID,
	}

	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to delete product biller: %w", err)
	}

	return nil
}

func (r *productBillerRepository) getBaseQuery(filters map[string]interface{}) (string, []interface{}) {
	var baseQuery = `
		SELECT id, product_id, biller_id, is_active, created_at, created_by, updated_at, updated_by
		FROM product_billers
	`

	var conditions []string
	var args []interface{}

	conditions = append(conditions, "deleted_at IS NULL")

	if productID, ok := filters["product_id"].(uint); ok {
		conditions = append(conditions, "product_id = ?")
		args = append(args, productID)
	}
	if billerID, ok := filters["biller_id"].(uint); ok {
		conditions = append(conditions, "biller_id = ?")
		args = append(args, billerID)
	}
	if isActive, ok := filters["is_active"].(bool); ok {
		conditions = append(conditions, "is_active = ?")
		args = append(args, isActive)
	}

	if len(conditions) > 0 {
		baseQuery = fmt.Sprintf("%s WHERE %s", baseQuery, strings.Join(conditions, " AND "))
	}

	return baseQuery, args
}

func (r *productBillerRepository) GetOne(ctx context.Context, id uint) (*models.ProductBiller, error) {
	query, args := r.getBaseQuery(map[string]interface{}{
		"id": id,
	})

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
