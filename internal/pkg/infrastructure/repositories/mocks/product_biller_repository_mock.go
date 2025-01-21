package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/models"
)

type MockProductBillerRepository struct {
	mock.Mock
}

func (m *MockProductBillerRepository) Create(ctx context.Context, productBiller *models.ProductBiller) error {
	args := m.Called(ctx, productBiller)
	return args.Error(0)
}

func (m *MockProductBillerRepository) Update(ctx context.Context, id int, productBiller *models.ProductBiller) error {
	args := m.Called(ctx, id, productBiller)
	return args.Error(0)
}

func (m *MockProductBillerRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductBillerRepository) DeleteByProductID(ctx context.Context, productID int) error {
	args := m.Called(ctx, productID)
	return args.Error(0)
}

func (m *MockProductBillerRepository) DeleteByBillerID(ctx context.Context, billerID int) error {
	args := m.Called(ctx, billerID)
	return args.Error(0)
}

func (m *MockProductBillerRepository) FetchOne(ctx context.Context, id int) (*models.ProductBiller, error) {
	args := m.Called(ctx, id)
	if pb, ok := args.Get(0).(*models.ProductBiller); ok {
		return pb, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductBillerRepository) FetchMany(ctx context.Context, filter map[string]interface{}) ([]*models.ProductBiller, error) {
	args := m.Called(ctx, filter)
	if pbs, ok := args.Get(0).([]*models.ProductBiller); ok {
		return pbs, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductBillerRepository) FetchManyWithPagination(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.ProductBiller, *db.Pagination, error) {
	args := m.Called(ctx, filter, page, limit)
	if pbs, ok := args.Get(0).([]*models.ProductBiller); ok {
		if pagination, ok := args.Get(1).(*db.Pagination); ok {
			return pbs, pagination, args.Error(2)
		}
	}
	return nil, nil, args.Error(2)
}
