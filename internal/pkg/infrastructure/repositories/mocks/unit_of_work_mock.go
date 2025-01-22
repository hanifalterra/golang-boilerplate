package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"golang-boilerplate/internal/pkg/infrastructure/repositories"
)

type MockUnitOfWork struct {
	mock.Mock
}

func (m *MockUnitOfWork) Execute(ctx context.Context, fn func(uow repositories.UnitOfWork) error) error {
	args := m.Called(ctx, fn)
	if fn != nil {
		_ = fn(m) // Call the function with the mock
	}
	return args.Error(0)
}

func (m *MockUnitOfWork) ProductRepo() repositories.ProductRepository {
	args := m.Called()
	return args.Get(0).(repositories.ProductRepository)
}

func (m *MockUnitOfWork) BillerRepo() repositories.BillerRepository {
	args := m.Called()
	return args.Get(0).(repositories.BillerRepository)
}

func (m *MockUnitOfWork) ProductBillerRepo() repositories.ProductBillerRepository {
	args := m.Called()
	return args.Get(0).(repositories.ProductBillerRepository)
}
