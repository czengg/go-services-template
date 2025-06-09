package upwardli

import (
	"elevate/internal/logger"

	"go.uber.org/zap"
)

// export types
type SubscriptionTopic = subscriptionTopic
type Webhook = webhook

type Service struct {
	WebhookProcessor
}

func NewService(config Config, logger logger.Logger, repo Repository) *Service {

	partnerClient, err := NewPartnerClient(partnerClientConfig{
		Config: config,
	})
	if err != nil {
		logger.Fatal("failed to create partner client", zap.Error(err))
	}

	return &Service{
		WebhookProcessor: NewWebhookProcessor(logger, partnerClient, repo),
	}
}
