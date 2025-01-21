package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/models"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) Update(ctx context.Context, id int, product *models.Product) error {
	args := m.Called(ctx, id, product)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepository) FetchOne(ctx context.Context, id int) (*models.Product, error) {
	args := m.Called(ctx, id)
	if p, ok := args.Get(0).(*models.Product); ok {
		return p, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductRepository) FetchMany(ctx context.Context, filter map[string]interface{}) ([]*models.Product, error) {
	args := m.Called(ctx, filter)
	if p, ok := args.Get(0).([]*models.Product); ok {
		return p, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductRepository) FetchManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Product, *db.Pagination, error) {
	args := m.Called(ctx, filter, page, limit)
	if p, ok := args.Get(0).([]*models.Product); ok {
		if pagination, ok := args.Get(1).(*db.Pagination); ok {
			return p, pagination, args.Error(2)
		}
	}
	return nil, nil, args.Error(2)
}
