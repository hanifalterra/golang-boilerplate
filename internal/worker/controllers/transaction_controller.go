package controllers

import (
	"context"
	"encoding/json"
	"fmt"

	"golang-boilerplate/internal/pkg/models"
	"golang-boilerplate/internal/pkg/services"
)

// TransactionController handles Kafka messages and parses them into Transaction structs.
type TransactionController struct {
	service *services.TransactionService
}

// NewTransactionController creates a new TransactionController.
func NewTransactionController(service *services.TransactionService) *TransactionController {
	return &TransactionController{service: service}
}

// HandleMessage parses a Kafka message and sends the Transaction to the service layer.
func (tc *TransactionController) HandleMessage(ctx context.Context, message []byte) error {
	var transaction models.Transaction
	if err := json.Unmarshal(message, &transaction); err != nil {
		return fmt.Errorf("failed to parse message: %w", err)
	}

	return tc.service.ProcessTransaction(ctx, transaction)
}
