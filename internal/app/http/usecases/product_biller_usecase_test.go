package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"golang-boilerplate/internal/app/http/usecases"
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
