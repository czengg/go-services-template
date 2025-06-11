package repository

import (
	"context"
	"database/sql"
	"template/internal/adapters/outbound/persistence/mysql/sqlc"
	"template/internal/core/upwardli"
	"template/packages/common-go"
)

func (r *repository) CreateUpwardliWebhook(ctx context.Context, webhook upwardli.Webhook) error {
	return r.queries.CreateUpwardliWebhook(ctx, sqlc.CreateUpwardliWebhookParams{
		ID:          webhook.ID,
		WebhookName: string(webhook.WebhookName),
		Endpoint:    webhook.Endpoint,
		PartnerID:   webhook.PartnerID,
		Status:      webhook.Status,
		Failures:    sql.NullInt32{Int32: int32(webhook.Failures), Valid: webhook.Failures != 0},
		LastFailure: sql.NullTime{Time: common.TimePtrToTime(webhook.LastFailure), Valid: webhook.LastFailure != nil},
	})
}

func (r *repository) GetAllUpwardliWebhooks(ctx context.Context) ([]upwardli.Webhook, error) {
	var webhooks []upwardli.Webhook
	rows, err := r.queries.GetAllUpwardliWebhooks(ctx)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		webhooks = append(webhooks, upwardli.Webhook{
			ID:          row.ID,
			WebhookName: upwardli.SubscriptionTopic(row.WebhookName),
			Endpoint:    row.Endpoint,
			PartnerID:   row.PartnerID,
			Status:      row.Status,
			Failures:    int64(row.Failures.Int32),
			LastFailure: common.TimeToTimePtr(row.LastFailure.Time),
		})
	}

	return webhooks, nil
}

func (r *repository) SoftDeleteUpwardliWebhook(ctx context.Context, id string) error {
	return r.queries.SoftDeleteUpwardliWebhook(ctx, id)
}
