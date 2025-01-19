package notification

import (
	"context"

	"golang-boilerplate/internal/pkg/connections/cacabot"
	"golang-boilerplate/internal/pkg/models"
)

type Notification interface {
	SendProductBillerSummary(ctx context.Context, payload models.ProductBillerSummaryNotification) error
}

type notification struct {
	client *cacabot.Client
}

func NewNotification(client *cacabot.Client) Notification {
	return &notification{client: client}
}

func (n *notification) SendProductBillerSummary(ctx context.Context, payload models.ProductBillerSummaryNotification) error {
	return n.client.SendMessage(ctx, "/product-biller-summary", payload)
}
