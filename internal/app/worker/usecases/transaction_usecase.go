package usecases

import (
	"context"
	"fmt"

	"golang-boilerplate/internal/pkg/infrastructure/lock"
	"golang-boilerplate/internal/pkg/infrastructure/repositories"
	"golang-boilerplate/internal/pkg/logger"
	"golang-boilerplate/internal/pkg/models"
)

type TransactionUseCase interface {
	ProcessTransaction(ctx context.Context, transaction *models.Transaction) error
}

type transactionUseCase struct {
	pbRepo repositories.ProductBillerRepository
	lock   lock.Lock
}

func NewTransactionUseCase(pbRepo repositories.ProductBillerRepository, lock lock.Lock) TransactionUseCase {
	return &transactionUseCase{
		pbRepo: pbRepo,
		lock:   lock,
	}
}

const eventClassTransaction = "usecase.transaction"

// ProcessTransaction processes a transaction.
func (uc *transactionUseCase) ProcessTransaction(ctx context.Context, transaction *models.Transaction) error {
	if transaction.Status == "success" {
		return nil
	}

	// Generate the lock key
	lockKey := fmt.Sprintf("worker:transaction:process_transaction:%d:%d", transaction.ProductID, transaction.BillerID)

	// Attempt to acquire the lock
	acquired, err := uc.lock.AcquireLock(ctx, lockKey)
	if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}
	if !acquired {
		return fmt.Errorf("failed to acquire lock within max time: lockKey=%q", lockKey)
	}

	// Ensure lock is released
	defer func() {
		if err := uc.lock.ReleaseLock(ctx, lockKey); err != nil {
			logger.FromContext(ctx).Error(ctx, eventClassTransaction, "ProcessTransaction.ReleaseLock", "[LockKey: %q]: %s", lockKey, err.Error())
		}
	}()

	// Fetch product-biller data
	pbs, err := uc.pbRepo.FetchMany(ctx, map[string]interface{}{
		"product_id": transaction.ProductID,
		"biller_id":  transaction.BillerID,
	})
	if err != nil {
		return fmt.Errorf("failed to fetch product-biller data: %w", err)
	}
	if len(pbs) == 0 {
		return fmt.Errorf("no product-biller data found")
	}

	// Check if the product-biller is active
	pb := pbs[0]
	if !pb.IsActive {
		return nil
	}

	// Deactivate the product-biller
	if err := uc.pbRepo.Update(ctx, pb.ID, &models.ProductBiller{IsActive: false}); err != nil {
		return fmt.Errorf("failed to deactivate product-biller: %w", err)
	}

	return nil
}
