package usecases

import (
	"context"
	"fmt"
	"time"

	"golang-boilerplate/internal/pkg/infrastructure/notification"
	"golang-boilerplate/internal/pkg/infrastructure/repositories"
	"golang-boilerplate/internal/pkg/models"
)

type CronUseCase struct {
	pbRepo repositories.ProductBillerRepository
	notif  notification.Notification
}

func NewCronUseCase(pbRepo repositories.ProductBillerRepository, notif notification.Notification) *CronUseCase {
	return &CronUseCase{
		pbRepo: pbRepo,
		notif:  notif,
	}
}

func (uc *CronUseCase) NotifyProductBillerSummary(ctx context.Context) error {
	// Fetch all product billers
	productBillers, err := uc.pbRepo.FetchMany(ctx, map[string]interface{}{})
	if err != nil {
		return fmt.Errorf("failed to fetch product billers: %w", err)
	}

	// Initialize the summary with current timestamp
	summary := models.ProductBillerSummaryNotification{
		DateTime: time.Now(),
		Total:    len(productBillers),
	}

	// Calculate active and inactive counts
	for _, productBiller := range productBillers {
		if productBiller.IsActive {
			summary.Active++
		}
	}

	// Inactive count is derived directly for efficiency
	summary.Inactive = summary.Total - summary.Active

	// Send the summary notification
	if err := uc.notif.SendProductBillerSummary(ctx, summary); err != nil {
		return fmt.Errorf("failed to send product biller summary notification: %w", err)
	}

	return nil
}
