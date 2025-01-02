package service

import (
	"context"
	"errors"

	db "golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/infrastructure/repository"
	"golang-boilerplate/internal/pkg/models"
)

// ProductBillerService defines the interface for the service layer of ProductBiller entities.
type ProductBillerService interface {
	Create(ctx context.Context, productBiller *models.ProductBiller) error
	Update(ctx context.Context, productID, billerID uint, productBiller *models.ProductBiller) error
	Delete(ctx context.Context, productID, billerID uint) error
	GetOne(ctx context.Context, productID, billerID uint) (*models.ProductBiller, error)
	GetMany(ctx context.Context, filter map[string]interface{}) ([]*models.ProductBiller, error)
	GetManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.ProductBiller, *db.Pagination, error)
}

// productBillerService implements ProductBillerService.
type productBillerService struct {
	repo repository.ProductBillerRepository
}

// NewProductBillerService creates a new instance of ProductBillerService.
func NewProductBillerService(repo repository.ProductBillerRepository) ProductBillerService {
	return &productBillerService{
		repo: repo,
	}
}

func (s *productBillerService) Create(ctx context.Context, productBiller *models.ProductBiller) error {
	if productBiller == nil {
		return errors.New("product biller is nil")
	}

	return s.repo.Create(ctx, productBiller)
}

func (s *productBillerService) Update(ctx context.Context, productID, billerID uint, productBiller *models.ProductBiller) error {
	if productBiller == nil {
		return errors.New("product biller is nil")
	}

	return s.repo.Update(ctx, productID, billerID, productBiller)
}

func (s *productBillerService) Delete(ctx context.Context, productID, billerID uint) error {
	return s.repo.Delete(ctx, productID, billerID)
}

func (s *productBillerService) GetOne(ctx context.Context, productID, billerID uint) (*models.ProductBiller, error) {
	return s.repo.GetOne(ctx, productID, billerID)
}

func (s *productBillerService) GetMany(ctx context.Context, filter map[string]interface{}) ([]*models.ProductBiller, error) {
	return s.repo.GetMany(ctx, filter)
}

func (s *productBillerService) GetManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.ProductBiller, *db.Pagination, error) {
	return s.repo.GetManyWithPagination(ctx, filter, page, limit)
}
