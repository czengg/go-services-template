package upwardli

import (
	"context"
)

// external types
type Repository interface {
	GetAllUpwardliWebhooks(ctx context.Context) ([]Webhook, error)
	CreateUpwardliWebhook(ctx context.Context, webhook Webhook) error
	SoftDeleteUpwardliWebhook(ctx context.Context, id string) error
}

type PartnerClient interface {
	GetAllWebhooks(ctx context.Context) ([]Webhook, error)
	CreateWebhook(ctx context.Context, endpoint string, topic string) (*Webhook, error)
	DeleteWebhook(ctx context.Context, webhookID string) error
}

type SubscriptionTopic = subscriptionTopic
type Webhook = webhook
