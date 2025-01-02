package repository

import (
	"context"
	"fmt"

	db "golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/models"
	utils_db "golang-boilerplate/internal/pkg/utils/db"
)

// BillerRepository defines the interface for managing Biller entities.
type BillerRepository interface {
	Create(ctx context.Context, biller *models.Biller) error
	Update(ctx context.Context, id uint, biller *models.Biller) error
	Delete(ctx context.Context, id uint) error
	GetOne(ctx context.Context, id uint) (*models.Biller, error)
	GetMany(ctx context.Context, filter map[string]interface{}) ([]*models.Biller, error)
	GetManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Biller, *db.Pagination, error)
}

// billerRepository implements BillerRepository.
type billerRepository struct {
	db db.DBExecutor
}

// NewBillerRepository creates a new instance of BillerRepository.
func NewBillerRepository(db db.DBExecutor) BillerRepository {
	return &billerRepository{
		db: db,
	}
}

func (r *billerRepository) Create(ctx context.Context, biller *models.Biller) error {
	const query = `
		INSERT INTO billers 
		(label, created_at, created_by, updated_at, updated_by)
		VALUES (:label, NOW(6), :created_by, NOW(6), :updated_by)
	`

	_, err := r.db.NamedExecContext(ctx, query, biller)
	if err != nil {
		return fmt.Errorf("failed to create biller: %w", err)
	}

	return nil
}

func (r *billerRepository) Update(ctx context.Context, id uint, biller *models.Biller) error {
	const query = `
		UPDATE billers
		SET label = :label, updated_at = NOW(6)
		WHERE id = :id AND deleted_at IS NULL
	`

	params := map[string]interface{}{
		"id":    id,
		"label": biller.Label,
	}

	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to update biller: %w", err)
	}

	return nil
}

func (r *billerRepository) Delete(ctx context.Context, id uint) error {
	const query = `
		UPDATE billers
		SET deleted_at = NOW(6)
		WHERE id = :id AND deleted_at IS NULL
	`

	params := map[string]interface{}{
		"id": id,
	}

	_, err := r.db.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("failed to delete biller: %w", err)
	}

	return nil
}

func (r *billerRepository) getBaseQuery(filter map[string]interface{}) (string, []interface{}) {
	const baseQuery = `
		SELECT product_id, biller_id, is_active, created_at, created_by, updated_at, updated_by
		FROM billers
		WHERE deleted_at IS NULL
	`

	return utils_db.ApplyFilters(baseQuery, filter)
}

func (r *billerRepository) GetOne(ctx context.Context, id uint) (*models.Biller, error) {
	filter := map[string]interface{}{
		"id": id,
	}

	query, args := r.getBaseQuery(filter)

	var biller models.Biller
	if err := r.db.GetContext(ctx, &biller, query, args...); err != nil {
		return nil, fmt.Errorf("failed to fetch biller: %w", err)
	}

	return &biller, nil
}

func (r *billerRepository) GetMany(ctx context.Context, filter map[string]interface{}) ([]*models.Biller, error) {
	query, args := r.getBaseQuery(filter)

	var billers []*models.Biller
	if err := r.db.SelectContext(ctx, &billers, query, args...); err != nil {
		return nil, fmt.Errorf("failed to fetch billers: %w", err)
	}

	return billers, nil
}

func (r *billerRepository) GetManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Biller, *db.Pagination, error) {
	query, args := r.getBaseQuery(filter)

	pagination := &db.Pagination{Order: "id ASC", Page: page, Limit: limit}
	var billers []*models.Biller
	if err := db.Paginate(ctx, r.db, query, args, pagination, &billers); err != nil {
		return nil, nil, fmt.Errorf("failed to fetch billers with pagination: %w", err)
	}

	return billers, pagination, nil
}
