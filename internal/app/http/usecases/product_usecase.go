package usecases

import (
	"context"
	"errors"
	"fmt"

	"golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/infrastructure/repositories"
	"golang-boilerplate/internal/pkg/models"
)

// ProductUseCase defines the interface for the usecase layer of Product entities.
type ProductUseCase interface {
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, id int, product *models.Product) error
	Delete(ctx context.Context, id int) error
	FetchOne(ctx context.Context, id int) (*models.Product, error)
	FetchMany(ctx context.Context, filter map[string]interface{}) ([]*models.Product, error)
	FetchManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Product, *db.Pagination, error)
}

// productUseCase implements ProductUseCase.
type productUseCase struct {
	repo repositories.ProductRepository
	uow  repositories.UnitOfWork
}

// NewProductUseCase creates a new instance of ProductUseCase.
func NewProductUseCase(repo repositories.ProductRepository, uow repositories.UnitOfWork) ProductUseCase {
	return &productUseCase{
		repo: repo,
		uow:  uow,
	}
}

func (uc *productUseCase) Create(ctx context.Context, product *models.Product) error {
	if product == nil {
		return errors.New("product is nil")
	}

	return uc.repo.Create(ctx, product)
}

func (uc *productUseCase) Update(ctx context.Context, id int, product *models.Product) error {
	if product == nil {
		return errors.New("product is nil")
	}

	return uc.repo.Update(ctx, id, product)
}

func (uc *productUseCase) Delete(ctx context.Context, id int) error {
	err := uc.uow.Execute(ctx, func(uow repositories.UnitOfWork) error {
		if err := uow.ProductBillerRepo().DeleteByProductID(ctx, id); err != nil {
			return fmt.Errorf("failed to delete product billers for product ID %d: %w", id, err)
		}

		if err := uow.ProductRepo().Delete(ctx, id); err != nil {
			return fmt.Errorf("failed to delete product with ID %d: %w", id, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("transaction failed while deleting product with ID %d: %w", id, err)
	}

	return nil
}

func (uc *productUseCase) FetchOne(ctx context.Context, id int) (*models.Product, error) {
	return uc.repo.FetchOne(ctx, id)
}

func (uc *productUseCase) FetchMany(ctx context.Context, filter map[string]interface{}) ([]*models.Product, error) {
	return uc.repo.FetchMany(ctx, filter)
}

func (uc *productUseCase) FetchManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Product, *db.Pagination, error) {
	return uc.repo.FetchManyWithPagination(ctx, filter, page, limit)
}
