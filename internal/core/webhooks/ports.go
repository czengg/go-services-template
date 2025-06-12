package webhooks

import (
	"context"
)

// external types
type Processor interface {
	Process(ctx context.Context, body []byte, headers map[string]string) error
}

type Verifier interface {
	Verify(body []byte, headers map[string]string) error
}

type Repository interface {
	GetAllWebhooksByProvider(ctx context.Context, provider provider) ([]Webhook, error)
	CreateWebhook(ctx context.Context, webhook Webhook) error
	SoftDeleteWebhook(ctx context.Context, provider Provider, id string) error
}

type SubscriptionClient interface {
	GetAllWebhooks(ctx context.Context) ([]Webhook, error)
	CreateWebhook(ctx context.Context, endpoint string, topic string) (*Webhook, error)
	DeleteWebhook(ctx context.Context, webhookID string) error
}

type WebhookManager interface {
	CreateWebhooks(ctx context.Context, endpoint string, topics []SubscriptionTopic) error
	CreateWebhook(ctx context.Context, endpoint string, topicName SubscriptionTopic) error
	GetWebhooks(ctx context.Context) ([]Webhook, error)
	DeleteWebhook(ctx context.Context, id string) error
}

type Service interface {
	WebhookManager
}

type SubscriptionTopic = subscriptionTopic
type Webhook = webhook
type Provider = provider

const (
	ProviderApril    provider = "april"
	ProviderUpwardli provider = "upwardli"
)
