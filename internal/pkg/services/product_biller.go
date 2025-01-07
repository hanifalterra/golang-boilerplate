package service

import (
	"context"
	"errors"
	"fmt"

	db "golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/infrastructure/repository"
	"golang-boilerplate/internal/pkg/models"
)

// ProductBillerService defines the interface for the service layer of ProductBiller entities.
type ProductBillerService interface {
	Create(ctx context.Context, productBiller *models.ProductBiller) error
	Update(ctx context.Context, id uint, productBiller *models.ProductBiller) error
	Delete(ctx context.Context, id uint) error
	GetOne(ctx context.Context, id uint) (*models.ProductBiller, error)
	GetMany(ctx context.Context, filter map[string]interface{}) ([]*models.ProductBiller, error)
	GetManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.ProductBiller, *db.Pagination, error)
}

// productBillerService implements ProductBillerService.
type productBillerService struct {
	repo        repository.ProductBillerRepository
	productRepo repository.ProductRepository
	billerRepo  repository.BillerRepository
}

// NewProductBillerService creates a new instance of ProductBillerService.
func NewProductBillerService(
	repo repository.ProductBillerRepository,
	productRepo repository.ProductRepository,
	billerRepo repository.BillerRepository,
) ProductBillerService {
	return &productBillerService{
		repo:        repo,
		productRepo: productRepo,
		billerRepo:  billerRepo,
	}
}

func (s *productBillerService) Create(ctx context.Context, productBiller *models.ProductBiller) error {
	if productBiller == nil {
		return errors.New("product biller is nil")
	}

	_, err := s.productRepo.GetOne(ctx, productBiller.ProductID)
	if err != nil {
		return fmt.Errorf("failed to fetch product with ID %d: %w", productBiller.ProductID, err)
	}

	_, err = s.billerRepo.GetOne(ctx, productBiller.BillerID)
	if err != nil {
		return fmt.Errorf("failed to fetch biller with ID %d: %w", productBiller.BillerID, err)
	}

	if err := s.repo.Create(ctx, productBiller); err != nil {
		return fmt.Errorf("failed to create product biller: %w", err)
	}

	return nil
}

func (s *productBillerService) Update(ctx context.Context, id uint, productBiller *models.ProductBiller) error {
	if productBiller == nil {
		return errors.New("product biller is nil")
	}

	return s.repo.Update(ctx, id, productBiller)
}

func (s *productBillerService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *productBillerService) GetOne(ctx context.Context, id uint) (*models.ProductBiller, error) {
	return s.repo.GetOne(ctx, id)
}

func (s *productBillerService) GetMany(ctx context.Context, filter map[string]interface{}) ([]*models.ProductBiller, error) {
	return s.repo.GetMany(ctx, filter)
}

func (s *productBillerService) GetManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.ProductBiller, *db.Pagination, error) {
	return s.repo.GetManyWithPagination(ctx, filter, page, limit)
}
