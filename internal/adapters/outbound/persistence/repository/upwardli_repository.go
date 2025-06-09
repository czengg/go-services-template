package repository

import (
	"context"
	"template/internal/core/upwardli"
)

func (r *repository) CreateUpwardliWebhook(ctx context.Context, webhook upwardli.Webhook) error {
	return nil
}

func (r *repository) GetAllUpwardliWebhooks(ctx context.Context) ([]upwardli.Webhook, error) {
	return nil, nil
}

func (r *repository) SoftDeleteUpwardliWebhook(ctx context.Context, id string) error {
	return nil
}
