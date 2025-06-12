package webhooks

import (
	"context"
	"strings"
	"template/internal/logger"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type webhookManager struct {
	logger   logger.Logger
	client   SubscriptionClient
	repo     Repository
	provider provider
}

func NewWebhookManager(
	logger logger.Logger,
	client SubscriptionClient,
	repo Repository,
	provider provider,
) WebhookManager {
	if logger == nil {
		return nil
	}

	return &webhookManager{
		logger:   logger,
		client:   client,
		repo:     repo,
		provider: provider,
	}
}

func (w *webhookManager) CreateWebhooks(ctx context.Context, endpoint string, topics []SubscriptionTopic) error {
	var errs []string
	successCount := 0

	for _, topic := range topics {
		err := w.CreateWebhook(ctx, endpoint, topic)
		if err != nil {
			w.logger.Error("failed to create webhook",
				zap.Error(err),
				zap.String("topic", string(topic)))
			errs = append(errs, string(topic))
		} else {
			successCount++
		}
	}

	if len(errs) > 0 {
		return errors.Errorf("failed to create %d webhooks for topics: %s",
			len(errs), strings.Join(errs, ", "))
	}

	w.logger.Info("successfully created all webhooks",
		zap.String("endpoint", endpoint),
		zap.Int("count", successCount))

	return nil
}

func (w *webhookManager) CreateWebhook(ctx context.Context, endpoint string, topicName SubscriptionTopic) error {
	if endpoint == "" {
		return errors.New("endpoint is required")
	}

	resp, err := w.client.CreateWebhook(ctx, endpoint, string(topicName))
	if err != nil {
		return errors.Wrap(err, "failed to create webhook via Upwardli")
	}

	// Get all webhooks to find our newly created one
	webhooksFromClient, err := w.client.GetAllWebhooks(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get webhooks from Upwardli")
	}

	// Find our webhook and save to database
	var webhookToSave *Webhook
	for _, webhook := range webhooksFromClient {
		if webhook.ID == resp.RegistrationID {
			webhookToSave = &Webhook{
				ID:          webhook.ID,
				WebhookName: webhook.WebhookName,
				Endpoint:    webhook.Endpoint,
				PartnerID:   webhook.PartnerID,
				Status:      webhook.Status,
				Failures:    webhook.Failures,
				LastFailure: webhook.LastFailure,
			}
			break
		}
	}

	if webhookToSave == nil {
		return errors.Errorf("webhook with registration ID %s not found in response", resp.RegistrationID)
	}

	// Save to database
	err = w.repo.CreateWebhook(ctx, *webhookToSave)
	if err != nil {
		return errors.Wrap(err, "failed to save webhook to database")
	}

	w.logger.Info("successfully created webhook",
		zap.String("topic", string(topicName)),
		zap.String("endpoint", endpoint),
		zap.String("webhookID", resp.RegistrationID))

	return nil
}

func (w *webhookManager) GetWebhooks(ctx context.Context) ([]Webhook, error) {
	webhooks, err := w.repo.GetAllWebhooksByProvider(ctx, w.provider)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get webhooks from database")
	}

	return webhooks, nil
}

func (w *webhookManager) DeleteWebhook(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("webhook ID is required")
	}

	// Delete from Upwardli first
	err := w.client.DeleteWebhook(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete webhook from Upwardli")
	}

	// Soft delete from database
	err = w.repo.SoftDeleteWebhook(ctx, w.provider, id)
	if err != nil {
		w.logger.Error("failed to delete webhook from database, but deleted from Upwardli",
			zap.Error(err),
			zap.String("webhookID", id),
			zap.String("provider", string(w.provider)))
		return errors.Wrap(err, "failed to delete webhook from database")
	}

	w.logger.Info("successfully deleted webhook", zap.String("webhookID", id))
	return nil
}
