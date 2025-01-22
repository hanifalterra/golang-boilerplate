package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"golang-boilerplate/internal/app/http/usecases"
	"golang-boilerplate/internal/pkg/connections/db"
	"golang-boilerplate/internal/pkg/infrastructure/repositories/mocks"
	"golang-boilerplate/internal/pkg/models"
)

func TestProductBillerUseCase_Create(t *testing.T) {
	type args struct {
		ctx           context.Context
		productBiller *models.ProductBiller
	}
	tests := []struct {
		name        string
		args        args
		setupMocks  func(productRepo *mocks.MockProductRepository, billerRepo *mocks.MockBillerRepository, repo *mocks.MockProductBillerRepository)
		expectErr   bool
		expectedErr error
	}{
		{
			name: "success",
			args: args{
				ctx: context.TODO(),
				productBiller: &models.ProductBiller{
					ProductID: 1,
					BillerID:  1,
				},
			},
			setupMocks: func(productRepo *mocks.MockProductRepository, billerRepo *mocks.MockBillerRepository, repo *mocks.MockProductBillerRepository) {
				productRepo.On("FetchOne", mock.Anything, 1).Return(&models.Product{}, nil)
				billerRepo.On("FetchOne", mock.Anything, 1).Return(&models.Biller{}, nil)
				repo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectErr:   false,
			expectedErr: nil,
		},
		{
			name: "error fetching product",
			args: args{
				ctx: context.TODO(),
				productBiller: &models.ProductBiller{
					ProductID: 1,
					BillerID:  1,
				},
			},
			setupMocks: func(productRepo *mocks.MockProductRepository, billerRepo *mocks.MockBillerRepository, repo *mocks.MockProductBillerRepository) {
				productRepo.On("FetchOne", mock.Anything, 1).Return(nil, errors.New("product not found"))
			},
			expectErr:   true,
			expectedErr: errors.New("failed to fetch product with ID 1: product not found"),
		},
		{
			name: "error fetching biller",
			args: args{
				ctx: context.TODO(),
				productBiller: &models.ProductBiller{
					ProductID: 1,
					BillerID:  1,
				},
			},
			setupMocks: func(productRepo *mocks.MockProductRepository, billerRepo *mocks.MockBillerRepository, repo *mocks.MockProductBillerRepository) {
				productRepo.On("FetchOne", mock.Anything, 1).Return(&models.Product{}, nil)
				billerRepo.On("FetchOne", mock.Anything, 1).Return(nil, errors.New("biller not found"))
			},
			expectErr:   true,
			expectedErr: errors.New("failed to fetch biller with ID 1: biller not found"),
		},
		{
			name: "error creating product biller",
			args: args{
				ctx: context.TODO(),
				productBiller: &models.ProductBiller{
					ProductID: 1,
					BillerID:  1,
				},
			},
			setupMocks: func(productRepo *mocks.MockProductRepository, billerRepo *mocks.MockBillerRepository, repo *mocks.MockProductBillerRepository) {
				productRepo.On("FetchOne", mock.Anything, 1).Return(&models.Product{}, nil)
				billerRepo.On("FetchOne", mock.Anything, 1).Return(&models.Biller{}, nil)
				repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("create failed"))
			},
			expectErr:   true,
			expectedErr: errors.New("failed to create product biller: create failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			productRepo := new(mocks.MockProductRepository)
			billerRepo := new(mocks.MockBillerRepository)
			repo := new(mocks.MockProductBillerRepository)

			uc := usecases.NewProductBillerUseCase(repo, productRepo, billerRepo)

			if tt.setupMocks != nil {
				tt.setupMocks(productRepo, billerRepo, repo)
			}

			err := uc.Create(tt.args.ctx, tt.args.productBiller)
			if tt.expectErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			productRepo.AssertExpectations(t)
			billerRepo.AssertExpectations(t)
			repo.AssertExpectations(t)
		})
	}
}

func TestProductBillerUseCase_Update(t *testing.T) {
	ctx := context.Background()

	id := 1
	updatedProductBiller := &models.ProductBiller{IsActive: true}

	t.Run("success", func(t *testing.T) {
		mockProductBillerRepo := new(mocks.MockProductBillerRepository)
		useCase := usecases.NewProductBillerUseCase(mockProductBillerRepo, nil, nil)

		mockProductBillerRepo.On("Update", ctx, id, updatedProductBiller).Return(nil)

		err := useCase.Update(ctx, id, updatedProductBiller)
		assert.NoError(t, err)

		mockProductBillerRepo.AssertCalled(t, "Update", ctx, id, updatedProductBiller)
	})

	t.Run("error", func(t *testing.T) {
		mockProductBillerRepo := new(mocks.MockProductBillerRepository)
		useCase := usecases.NewProductBillerUseCase(mockProductBillerRepo, nil, nil)

		mockProductBillerRepo.On("Update", ctx, id, updatedProductBiller).Return(errors.New("update failed"))

		err := useCase.Update(ctx, id, updatedProductBiller)
		assert.Error(t, err)

		mockProductBillerRepo.AssertCalled(t, "Update", ctx, id, updatedProductBiller)
	})

	t.Run("nil product biller", func(t *testing.T) {
		mockProductBillerRepo := new(mocks.MockProductBillerRepository)
		useCase := usecases.NewProductBillerUseCase(mockProductBillerRepo, nil, nil)

		err := useCase.Update(ctx, id, nil)
		assert.Error(t, err)
		assert.Equal(t, "product biller is nil", err.Error())

		mockProductBillerRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestProductBillerUseCase_Delete(t *testing.T) {
	ctx := context.Background()

	id := 1

	t.Run("success", func(t *testing.T) {
		mockProductBillerRepo := new(mocks.MockProductBillerRepository)
		useCase := usecases.NewProductBillerUseCase(mockProductBillerRepo, nil, nil)

		mockProductBillerRepo.On("Delete", ctx, id).Return(nil)

		err := useCase.Delete(ctx, id)
		assert.NoError(t, err)

		mockProductBillerRepo.AssertCalled(t, "Delete", ctx, id)
	})

	t.Run("error", func(t *testing.T) {
		mockProductBillerRepo := new(mocks.MockProductBillerRepository)
		useCase := usecases.NewProductBillerUseCase(mockProductBillerRepo, nil, nil)

		mockProductBillerRepo.On("Delete", ctx, id).Return(errors.New("delete failed"))

		err := useCase.Delete(ctx, id)
		assert.Error(t, err)

		mockProductBillerRepo.AssertCalled(t, "Delete", ctx, id)
	})
}

func TestProductBillerUseCase_FetchOne(t *testing.T) {
	ctx := context.Background()

	id := 1
	expected := &models.ProductBiller{ID: id, ProductID: 1, BillerID: 1, IsActive: true}

	t.Run("success", func(t *testing.T) {
		mockProductBillerRepo := new(mocks.MockProductBillerRepository)
		useCase := usecases.NewProductBillerUseCase(mockProductBillerRepo, nil, nil)

		mockProductBillerRepo.On("FetchOne", ctx, id).Return(expected, nil)

		result, err := useCase.FetchOne(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)

		mockProductBillerRepo.AssertCalled(t, "FetchOne", ctx, id)
	})

	t.Run("error", func(t *testing.T) {
		mockProductBillerRepo := new(mocks.MockProductBillerRepository)
		useCase := usecases.NewProductBillerUseCase(mockProductBillerRepo, nil, nil)

		mockProductBillerRepo.On("FetchOne", ctx, id).Return(nil, errors.New("not found"))

		result, err := useCase.FetchOne(ctx, id)
		assert.Error(t, err)
		assert.Nil(t, result)

		mockProductBillerRepo.AssertCalled(t, "FetchOne", ctx, id)
	})
}

func TestProductBillerUseCase_FetchMany(t *testing.T) {
	ctx := context.Background()

	filter := map[string]interface{}{"label": "Test"}
	expected := []*models.ProductBiller{
		{ID: 1, ProductID: 1, BillerID: 1, IsActive: true},
		{ID: 2, ProductID: 2, BillerID: 2, IsActive: false},
	}

	t.Run("success", func(t *testing.T) {
		mockProductBillerRepo := new(mocks.MockProductBillerRepository)
		useCase := usecases.NewProductBillerUseCase(mockProductBillerRepo, nil, nil)

		mockProductBillerRepo.On("FetchMany", ctx, filter).Return(expected, nil)

		result, err := useCase.FetchMany(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)

		mockProductBillerRepo.AssertCalled(t, "FetchMany", ctx, filter)
	})

	t.Run("error", func(t *testing.T) {
		mockProductBillerRepo := new(mocks.MockProductBillerRepository)
		useCase := usecases.NewProductBillerUseCase(mockProductBillerRepo, nil, nil)

		mockProductBillerRepo.On("FetchMany", ctx, filter).Return(nil, errors.New("fetch failed"))

		result, err := useCase.FetchMany(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, result)

		mockProductBillerRepo.AssertCalled(t, "FetchMany", ctx, filter)
	})
}

func TestProductBillerUseCase_FetchManyWithPagination(t *testing.T) {
	ctx := context.Background()

	filter := map[string]interface{}{"label": "Test"}
	page := 1
	limit := 10
	expectedData := []*models.ProductBiller{
		{ID: 1, ProductID: 1, BillerID: 1, IsActive: true},
		{ID: 2, ProductID: 2, BillerID: 2, IsActive: false},
	}
	expectedPagination := &db.Pagination{
		Limit: limit,
		Page:  page,
	}

	t.Run("success", func(t *testing.T) {
		mockProductBillerRepo := new(mocks.MockProductBillerRepository)
		useCase := usecases.NewProductBillerUseCase(mockProductBillerRepo, nil, nil)

		mockProductBillerRepo.On("FetchManyWithPagination", ctx, filter, page, limit).Return(expectedData, expectedPagination, nil)

		data, pagination, err := useCase.FetchManyWithPagination(ctx, filter, page, limit)
		assert.NoError(t, err)
		assert.Equal(t, expectedData, data)
		assert.Equal(t, expectedPagination, pagination)

		mockProductBillerRepo.AssertCalled(t, "FetchManyWithPagination", ctx, filter, page, limit)
	})

	t.Run("error", func(t *testing.T) {
		mockProductBillerRepo := new(mocks.MockProductBillerRepository)
		useCase := usecases.NewProductBillerUseCase(mockProductBillerRepo, nil, nil)

		mockProductBillerRepo.On("FetchManyWithPagination", ctx, filter, page, limit).Return(nil, nil, errors.New("fetch failed"))

		data, pagination, err := useCase.FetchManyWithPagination(ctx, filter, page, limit)
		assert.Error(t, err)
		assert.Nil(t, data)
		assert.Nil(t, pagination)

		mockProductBillerRepo.AssertCalled(t, "FetchManyWithPagination", ctx, filter, page, limit)
	})
}
