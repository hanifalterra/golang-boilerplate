package usecases

import (
	"context"
	"errors"
	"fmt"

	"golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/infrastructure/repositories"
	"golang-boilerplate/internal/pkg/models"
)

// ProductBillerUseCase defines the interface for the usecase layer of ProductBiller entities.
type ProductBillerUseCase interface {
	Create(ctx context.Context, productBiller *models.ProductBiller) error
	Update(ctx context.Context, id uint, productBiller *models.ProductBiller) error
	Delete(ctx context.Context, id uint) error
	FetchOne(ctx context.Context, id uint) (*models.ProductBiller, error)
	FetchMany(ctx context.Context, filter map[string]interface{}) ([]*models.ProductBiller, error)
	FetchManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.ProductBiller, *db.Pagination, error)
}

// productBillerUseCase implements ProductBillerUseCase.
type productBillerUseCase struct {
	repo        repositories.ProductBillerRepository
	productRepo repositories.ProductRepository
	billerRepo  repositories.BillerRepository
}

// NewProductBillerUseCase creates a new instance of ProductBillerUseCase.
func NewProductBillerUseCase(
	repo repositories.ProductBillerRepository,
	productRepo repositories.ProductRepository,
	billerRepo repositories.BillerRepository,
) ProductBillerUseCase {
	return &productBillerUseCase{
		repo:        repo,
		productRepo: productRepo,
		billerRepo:  billerRepo,
	}
}

func (uc *productBillerUseCase) Create(ctx context.Context, productBiller *models.ProductBiller) error {
	if productBiller == nil {
		return errors.New("product biller is nil")
	}

	_, err := uc.productRepo.FetchOne(ctx, productBiller.ProductID)
	if err != nil {
		return fmt.Errorf("failed to fetch product with ID %d: %w", productBiller.ProductID, err)
	}

	_, err = uc.billerRepo.FetchOne(ctx, productBiller.BillerID)
	if err != nil {
		return fmt.Errorf("failed to fetch biller with ID %d: %w", productBiller.BillerID, err)
	}

	if err := uc.repo.Create(ctx, productBiller); err != nil {
		return fmt.Errorf("failed to create product biller: %w", err)
	}

	return nil
}

func (uc *productBillerUseCase) Update(ctx context.Context, id uint, productBiller *models.ProductBiller) error {
	if productBiller == nil {
		return errors.New("product biller is nil")
	}

	return uc.repo.Update(ctx, id, productBiller)
}

func (uc *productBillerUseCase) Delete(ctx context.Context, id uint) error {
	return uc.repo.Delete(ctx, id)
}

func (uc *productBillerUseCase) FetchOne(ctx context.Context, id uint) (*models.ProductBiller, error) {
	return uc.repo.FetchOne(ctx, id)
}

func (uc *productBillerUseCase) FetchMany(ctx context.Context, filter map[string]interface{}) ([]*models.ProductBiller, error) {
	return uc.repo.FetchMany(ctx, filter)
}

func (uc *productBillerUseCase) FetchManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.ProductBiller, *db.Pagination, error) {
	return uc.repo.FetchManyWithPagination(ctx, filter, page, limit)
}
