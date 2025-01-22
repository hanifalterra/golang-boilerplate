package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockLock struct {
	mock.Mock
}

func (m *MockLock) AcquireLock(ctx context.Context, lockKey string) (bool, error) {
	args := m.Called(ctx, lockKey)
	return args.Bool(0), args.Error(1)
}

func (m *MockLock) ReleaseLock(ctx context.Context, lockKey string) error {
	args := m.Called(ctx, lockKey)
	return args.Error(0)
}
