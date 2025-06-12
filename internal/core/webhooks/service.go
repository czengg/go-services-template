package webhooks

import (
	"template/internal/logger"
)

type service struct {
	WebhookManager
}

func NewService(logger logger.Logger, repo Repository, client SubscriptionClient, provider provider) *service {
	return &service{
		WebhookManager: NewWebhookManager(logger, client, repo, provider),
	}
}
