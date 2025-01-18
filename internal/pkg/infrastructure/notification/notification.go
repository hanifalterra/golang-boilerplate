package notification

import (
	"context"

	"golang-boilerplate/internal/pkg/connections/cacabot"
	"golang-boilerplate/internal/pkg/models"
)

type Notification interface {
	SendCountProductBillers(ctx context.Context, payload models.CountProductBillerNotification) error
}

type notification struct {
	client *cacabot.Client
}

func NewNotification(client *cacabot.Client) Notification {
	return &notification{client: client}
}

func (n *notification) SendCountProductBillers(ctx context.Context, payload models.CountProductBillerNotification) error {
	return n.client.SendMessage(ctx, "/count-product-billers", payload)
}
