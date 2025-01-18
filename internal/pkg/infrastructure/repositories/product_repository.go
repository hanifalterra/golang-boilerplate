package repositories

import (
	"context"
	"fmt"
	"strings"

	"golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/models"
)

// ProductRepository defines the interface for managing Product entities.
type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, id uint, product *models.Product) error
	Delete(ctx context.Context, id uint) error
	FetchOne(ctx context.Context, id uint) (*models.Product, error)
	FetchMany(ctx context.Context, filter map[string]interface{}) ([]*models.Product, error)
	FetchManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Product, *db.Pagination, error)
}

// productRepository implements ProductRepository.
type productRepository struct {
	db db.DBExecutor
}

// NewProductRepository creates a new instance of ProductRepository.
func NewProductRepository(db db.DBExecutor) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	const query = `
		INSERT INTO products 
		(label, created_at, created_by, updated_at, updated_by)
		VALUES (:label, NOW(6), :created_by, NOW(6), :updated_by)
	`

	_, err := r.db.NamedExecContext(ctx, query, product)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

func (r *productRepository) Update(ctx context.Context, id uint, product *models.Product) error {
	const query = `
		UPDATE products
		SET label = :label, updated_at = NOW(6)
		WHERE id = :id AND deleted_at IS NULL
	`

	params := map[string]interface{}{
		"id":    id,
		"label": product.Label,
	}

	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

func (r *productRepository) Delete(ctx context.Context, id uint) error {
	const query = `
		UPDATE products
		SET deleted_at = NOW(6)
		WHERE id = :id AND deleted_at IS NULL
	`

	params := map[string]interface{}{
		"id": id,
	}

	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

func (r *productRepository) getBaseQuery(filters map[string]interface{}) (string, []interface{}) {
	var baseQuery = `
		SELECT id, label, created_at, created_by, updated_at, updated_by
		FROM products
	`

	var conditions []string
	var args []interface{}

	conditions = append(conditions, "deleted_at IS NULL")

	if id, ok := filters["id"].(uint); ok {
		conditions = append(conditions, "id = ?")
		args = append(args, id)
	}
	if label, ok := filters["label"].(string); ok {
		conditions = append(conditions, "label LIKE ?")
		args = append(args, "%"+label+"%")
	}

	if len(conditions) > 0 {
		baseQuery = fmt.Sprintf("%s WHERE %s", baseQuery, strings.Join(conditions, " AND "))
	}

	return baseQuery, args
}

func (r *productRepository) FetchOne(ctx context.Context, id uint) (*models.Product, error) {
	query, args := r.getBaseQuery(map[string]interface{}{
		"id": id,
	})

	var product models.Product
	if err := r.db.GetContext(ctx, &product, query, args...); err != nil {
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}

	return &product, nil
}

func (r *productRepository) FetchMany(ctx context.Context, filter map[string]interface{}) ([]*models.Product, error) {
	query, args := r.getBaseQuery(filter)

	var products []*models.Product
	if err := r.db.SelectContext(ctx, &products, query, args...); err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}

	return products, nil
}

func (r *productRepository) FetchManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Product, *db.Pagination, error) {
	query, args := r.getBaseQuery(filter)

	pagination := &db.Pagination{Order: "id ASC", Page: page, Limit: limit}
	var products []*models.Product
	if err := db.Paginate(ctx, r.db, query, args, pagination, &products); err != nil {
		return nil, nil, fmt.Errorf("failed to fetch products with pagination: %w", err)
	}

	return products, pagination, nil
}
