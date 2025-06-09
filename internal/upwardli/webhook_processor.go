package upwardli

import (
	"context"
	"elevate/internal/logger"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type WebhookProcessor interface {
	CreateAllWebhooks(ctx context.Context, endpoint string) error
	CreateWebhook(ctx context.Context, topicName SubscriptionTopic, endpoint string) error
	GetWebhooks(ctx context.Context) ([]Webhook, error)
	DeleteWebhook(ctx context.Context, id string) error
}

type webhookProcessor struct {
	logger        logger.Logger
	partnerClient PartnerClient
	repo          Repository
}

func NewWebhookProcessor(
	logger logger.Logger,
	partnerClient PartnerClient,
	repo Repository,
) WebhookProcessor {
	if logger == nil {
		return nil
	}

	return &webhookProcessor{
		logger:        logger,
		partnerClient: partnerClient,
		repo:          repo,
	}
}

func (w *webhookProcessor) CreateAllWebhooks(ctx context.Context, endpoint string) error {
	webhookTopics := []SubscriptionTopic{
		SubscriptionTopicConsumerCreated,
		SubscriptionTopicConsumerUpdated,
		SubscriptionTopicConsumerClosed,
		SubscriptionTopicConsumerKYCStarted,
		SubscriptionTopicConsumerKYCPending,
		SubscriptionTopicConsumerKYCCompleted,
		SubscriptionTopicConsumerKYCNeedsReview,
		SubscriptionTopicConsumerKYCApproved,
		SubscriptionTopicConsumerKYCFailed,
		SubscriptionTopicPaymentCardCreated,
		SubscriptionTopicPaymentCardUpdated,
		SubscriptionTopicPaymentCardClosed,
		SubscriptionTopicPaymentCardTransactionSettlement,
		SubscriptionTopicACHSent,
		SubscriptionTopicACHReceived,
		SubscriptionTopicACHFailed,
		SubscriptionTopicPaymentTransferCreated,
		SubscriptionTopicPaymentTransferCompleted,
		SubscriptionTopicPaymentTransferFailed,
	}

	var errs []string
	successCount := 0

	for _, topic := range webhookTopics {
		err := w.CreateWebhook(ctx, topic, endpoint)
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

func (w *webhookProcessor) CreateWebhook(ctx context.Context, topicName SubscriptionTopic, endpoint string) error {
	if endpoint == "" {
		return errors.New("endpoint is required")
	}

	// Create webhook via Upwardli API
	webhookReq := CreateWebhookRequest{
		WebhookName: topicName,
		Endpoint:    endpoint,
	}

	resp, err := w.partnerClient.CreateWebhook(ctx, webhookReq)
	if err != nil {
		return errors.Wrap(err, "failed to create webhook via Upwardli")
	}

	// Get all webhooks to find our newly created one
	webhooksFromClient, err := w.partnerClient.GetAllWebhooks(ctx)
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
	err = w.repo.CreateUpwardliWebhook(ctx, *webhookToSave)
	if err != nil {
		return errors.Wrap(err, "failed to save webhook to database")
	}

	w.logger.Info("successfully created webhook",
		zap.String("topic", string(topicName)),
		zap.String("endpoint", endpoint),
		zap.String("webhookID", resp.RegistrationID))

	return nil
}

func (w *webhookProcessor) GetWebhooks(ctx context.Context) ([]Webhook, error) {
	webhooks, err := w.repo.GetAllUpwardliWebhooks(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get webhooks from database")
	}

	return webhooks, nil
}

func (w *webhookProcessor) DeleteWebhook(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("webhook ID is required")
	}

	// Delete from Upwardli first
	err := w.partnerClient.DeleteWebhook(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete webhook from Upwardli")
	}

	// Soft delete from database
	err = w.repo.SoftDeleteUpwardliWebhook(ctx, id)
	if err != nil {
		w.logger.Error("failed to delete webhook from database, but deleted from Upwardli",
			zap.Error(err),
			zap.String("webhookID", id))
		return errors.Wrap(err, "failed to delete webhook from database")
	}

	w.logger.Info("successfully deleted webhook", zap.String("webhookID", id))
	return nil
}
