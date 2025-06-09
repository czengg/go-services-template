package upwardli

import (
	"context"
)

type Repository interface {
	GetAllUpwardliWebhooks(ctx context.Context) ([]Webhook, error)
	CreateUpwardliWebhook(ctx context.Context, webhook Webhook) error
	SoftDeleteUpwardliWebhook(ctx context.Context, id string) error
}
