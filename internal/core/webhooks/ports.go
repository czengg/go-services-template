package webhooks

import (
	"context"
)

// external types
type Repository interface {
	GetAllBankingWebhooks(ctx context.Context) ([]Webhook, error)
	CreateBankingWebhook(ctx context.Context, webhook Webhook) error
	SoftDeleteBankingWebhook(ctx context.Context, id string) error
}

type Client interface {
	GetAllWebhooks(ctx context.Context) ([]Webhook, error)
	CreateWebhook(ctx context.Context, endpoint string, topic string) (*Webhook, error)
	DeleteWebhook(ctx context.Context, webhookID string) error
}

type SubscriptionTopic = subscriptionTopic
type Webhook = webhook
