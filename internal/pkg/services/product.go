package service

import (
	"context"
	"errors"

	db "golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/infrastructure/repository"
	"golang-boilerplate/internal/pkg/models"
)

// ProductService defines the interface for the service layer of Product entities.
type ProductService interface {
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, id uint, product *models.Product) error
	Delete(ctx context.Context, id uint) error
	GetOne(ctx context.Context, id uint) (*models.Product, error)
	GetMany(ctx context.Context, filter map[string]interface{}) ([]*models.Product, error)
	GetManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Product, *db.Pagination, error)
}

// productService implements ProductService.
type productService struct {
	repo repository.ProductRepository
}

// NewProductService creates a new instance of ProductService.
func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) Create(ctx context.Context, product *models.Product) error {
	if product == nil {
		return errors.New("product is nil")
	}

	return s.repo.Create(ctx, product)
}

func (s *productService) Update(ctx context.Context, id uint, product *models.Product) error {
	if product == nil {
		return errors.New("product is nil")
	}

	return s.repo.Update(ctx, id, product)
}

func (s *productService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *productService) GetOne(ctx context.Context, id uint) (*models.Product, error) {
	return s.repo.GetOne(ctx, id)
}

func (s *productService) GetMany(ctx context.Context, filter map[string]interface{}) ([]*models.Product, error) {
	return s.repo.GetMany(ctx, filter)
}

func (s *productService) GetManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Product, *db.Pagination, error) {
	return s.repo.GetManyWithPagination(ctx, filter, page, limit)
}
