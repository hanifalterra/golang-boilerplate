package service

import (
	"context"
	"errors"

	db "golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/infrastructure/repository"
	"golang-boilerplate/internal/pkg/models"
)

// BillerService defines the interface for the service layer of Biller entities.
type BillerService interface {
	Create(ctx context.Context, biller *models.Biller) error
	Update(ctx context.Context, id uint, biller *models.Biller) error
	Delete(ctx context.Context, id uint) error
	GetOne(ctx context.Context, id uint) (*models.Biller, error)
	GetMany(ctx context.Context, filter map[string]interface{}) ([]*models.Biller, error)
	GetManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Biller, *db.Pagination, error)
}

// billerService implements BillerService.
type billerService struct {
	repo repository.BillerRepository
}

// NewBillerService creates a new instance of BillerService.
func NewBillerService(repo repository.BillerRepository) BillerService {
	return &billerService{
		repo: repo,
	}
}

func (s *billerService) Create(ctx context.Context, biller *models.Biller) error {
	if biller == nil {
		return errors.New("biller is nil")
	}

	return s.repo.Create(ctx, biller)
}

func (s *billerService) Update(ctx context.Context, id uint, biller *models.Biller) error {
	if biller == nil {
		return errors.New("biller is nil")
	}

	return s.repo.Update(ctx, id, biller)
}

func (s *billerService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func (s *billerService) GetOne(ctx context.Context, id uint) (*models.Biller, error) {
	return s.repo.GetOne(ctx, id)
}

func (s *billerService) GetMany(ctx context.Context, filter map[string]interface{}) ([]*models.Biller, error) {
	return s.repo.GetMany(ctx, filter)
}

func (s *billerService) GetManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Biller, *db.Pagination, error) {
	return s.repo.GetManyWithPagination(ctx, filter, page, limit)
}
