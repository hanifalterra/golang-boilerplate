package services

import (
	"context"
	"fmt"
	"golang-boilerplate/internal/pkg/infrastructure/notification"
	"golang-boilerplate/internal/pkg/infrastructure/repositories"
)

type ProductBillerService struct {
	Repository          repositories.ProductBillerRepository
	NotificationService notification.NotificationService
}

func NewProductBillerService(repo repositories.ProductBillerRepository, notif notification.NotificationService) *ProductBillerService {
	return &ProductBillerService{
		Repository:          repo,
		NotificationService: notif,
	}
}

func (s *ProductBillerService) ProcessInactiveBillers() error {
	billers, err := s.Repository.GetMany(context.Background(), map[string]interface{}{})
	if err != nil {
		return err
	}

	count := 0
	for _, biller := range billers {
		if !biller.IsActive {
			count++
		}
	}

	message := "Count of inactive ProductBillers: " + fmt.Sprint(count)
	return s.NotificationService.SendMessage(message)
}
