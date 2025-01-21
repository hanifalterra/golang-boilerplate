package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/models"
)

type MockBillerRepository struct {
	mock.Mock
}

func (m *MockBillerRepository) Create(ctx context.Context, biller *models.Biller) error {
	args := m.Called(ctx, biller)
	return args.Error(0)
}

func (m *MockBillerRepository) Update(ctx context.Context, id int, biller *models.Biller) error {
	args := m.Called(ctx, id, biller)
	return args.Error(0)
}

func (m *MockBillerRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockBillerRepository) FetchOne(ctx context.Context, id int) (*models.Biller, error) {
	args := m.Called(ctx, id)
	if b, ok := args.Get(0).(*models.Biller); ok {
		return b, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBillerRepository) FetchMany(ctx context.Context, filter map[string]interface{}) ([]*models.Biller, error) {
	args := m.Called(ctx, filter)
	if b, ok := args.Get(0).([]*models.Biller); ok {
		return b, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBillerRepository) FetchManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Biller, *db.Pagination, error) {
	args := m.Called(ctx, filter, page, limit)
	if b, ok := args.Get(0).([]*models.Biller); ok {
		if p, ok := args.Get(1).(*db.Pagination); ok {
			return b, p, args.Error(2)
		}
	}
	return nil, nil, args.Error(2)
}
