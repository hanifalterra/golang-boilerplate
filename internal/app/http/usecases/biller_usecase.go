package usecases

import (
	"context"
	"errors"
	"fmt"

	"golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/infrastructure/repositories"
	"golang-boilerplate/internal/pkg/models"
)

// BillerUseCase defines the interface for the usecase layer of Biller entities.
type BillerUseCase interface {
	Create(ctx context.Context, biller *models.Biller) error
	Update(ctx context.Context, id uint, biller *models.Biller) error
	Delete(ctx context.Context, id uint) error
	FetchOne(ctx context.Context, id uint) (*models.Biller, error)
	FetchMany(ctx context.Context, filter map[string]interface{}) ([]*models.Biller, error)
	FetchManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Biller, *db.Pagination, error)
}

// billerUseCase implements BillerUseCase.
type billerUseCase struct {
	repo repositories.BillerRepository
	uow  repositories.UnitOfWork
}

// NewBillerUseCase creates a new instance of BillerUseCase.
func NewBillerUseCase(repo repositories.BillerRepository, uow repositories.UnitOfWork) BillerUseCase {
	return &billerUseCase{
		repo: repo,
		uow:  uow,
	}
}

func (uc *billerUseCase) Create(ctx context.Context, biller *models.Biller) error {
	if biller == nil {
		return errors.New("biller is nil")
	}

	return uc.repo.Create(ctx, biller)
}

func (uc *billerUseCase) Update(ctx context.Context, id uint, biller *models.Biller) error {
	if biller == nil {
		return errors.New("biller is nil")
	}

	return uc.repo.Update(ctx, id, biller)
}

func (uc *billerUseCase) Delete(ctx context.Context, id uint) error {
	err := uc.uow.Execute(ctx, func(uow repositories.UnitOfWork) error {
		if err := uow.ProductBillerRepo().DeleteByBillerID(ctx, id); err != nil {
			return fmt.Errorf("failed to delete product billers for biller ID %d: %w", id, err)
		}

		if err := uow.BillerRepo().Delete(ctx, id); err != nil {
			return fmt.Errorf("failed to delete biller with ID %d: %w", id, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("transaction failed while deleting biller with ID %d: %w", id, err)
	}

	return nil
}

func (uc *billerUseCase) FetchOne(ctx context.Context, id uint) (*models.Biller, error) {
	return uc.repo.FetchOne(ctx, id)
}

func (uc *billerUseCase) FetchMany(ctx context.Context, filter map[string]interface{}) ([]*models.Biller, error) {
	return uc.repo.FetchMany(ctx, filter)
}

func (uc *billerUseCase) FetchManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Biller, *db.Pagination, error) {
	return uc.repo.FetchManyWithPagination(ctx, filter, page, limit)
}
