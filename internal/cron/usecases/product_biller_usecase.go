package usecases

import (
	"context"
	"golang-boilerplate/internal/pkg/infrastructure/notification"
	"golang-boilerplate/internal/pkg/infrastructure/repositories"
	"golang-boilerplate/internal/pkg/models"
	"time"
)

type ProductBillerUseCase struct {
	repo  repositories.ProductBillerRepository
	notif notification.Notification
}

func NewProductBillerUseCase(repo repositories.ProductBillerRepository, notif notification.Notification) *ProductBillerUseCase {
	return &ProductBillerUseCase{
		repo:  repo,
		notif: notif,
	}
}

func (uc *ProductBillerUseCase) ProcessCountProductBillers(ctx context.Context) error {
	pbs, err := uc.repo.FetchMany(ctx, map[string]interface{}{})
	if err != nil {
		return err
	}

	count := models.CountProductBillerNotification{
		DateTime: time.Now(),
	}
	for _, pb := range pbs {
		count.Total++
		if pb.IsActive {
			count.Active++
		} else {
			count.Inactive++
		}
	}

	return uc.notif.SendCountProductBillers(ctx, count)
}
