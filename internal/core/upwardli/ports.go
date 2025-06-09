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

type CreateWebhookRequest struct {
	WebhookName SubscriptionTopic `json:"webhook_name"`
	Endpoint    string            `json:"endpoint"`
}

type PartnerClient interface {
	GetAllWebhooks(ctx context.Context) ([]Webhook, error)
	CreateWebhook(ctx context.Context, webhook CreateWebhookRequest) (*Webhook, error)
	DeleteWebhook(ctx context.Context, webhookID string) error
}

type SubscriptionTopic = subscriptionTopic
type Webhook = webhook
