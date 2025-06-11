package repository

import (
	"context"
	"database/sql"
	"template/internal/adapters/outbound/persistence/mysql/sqlc"
	webhooks "template/internal/core/webhooks"
	"template/packages/common-go"
)

func (r *repository) CreateBankingWebhook(ctx context.Context, webhook webhooks.Webhook) error {
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

func (r *repository) GetAllBankingWebhooks(ctx context.Context) ([]webhooks.Webhook, error) {
	var ws []webhooks.Webhook
	rows, err := r.queries.GetAllUpwardliWebhooks(ctx)
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		ws = append(ws, webhooks.Webhook{
			ID:          row.ID,
			WebhookName: webhooks.SubscriptionTopic(row.WebhookName),
			Endpoint:    row.Endpoint,
			PartnerID:   row.PartnerID,
			Status:      row.Status,
			Failures:    int64(row.Failures.Int32),
			LastFailure: common.TimeToTimePtr(row.LastFailure.Time),
		})
	}

	return ws, nil
}

func (r *repository) SoftDeleteBankingWebhook(ctx context.Context, id string) error {
	return r.queries.SoftDeleteUpwardliWebhook(ctx, id)
}
