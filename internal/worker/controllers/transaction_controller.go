package controllers

import (
	"context"
	"encoding/json"
	"fmt"

	"golang-boilerplate/internal/pkg/models"
	"golang-boilerplate/internal/worker/usecases"
)

// TransactionController handles Kafka messages and parses them into Transaction structs.
type TransactionController struct {
	usecase usecases.TransactionUseCase
}

// NewTransactionController creates a new TransactionController.
func NewTransactionController(usecase usecases.TransactionUseCase) *TransactionController {
	return &TransactionController{usecase: usecase}
}

// HandleMessage parses a Kafka message and sends the Transaction to the usecase layer.
func (tc *TransactionController) HandleMessage(ctx context.Context, message []byte) error {
	var transaction *models.Transaction
	if err := json.Unmarshal(message, transaction); err != nil {
		return fmt.Errorf("failed to parse message: %w", err)
	}

	return tc.usecase.ProcessTransaction(ctx, transaction)
}
