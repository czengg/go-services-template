package webhooks

import (
	"template/internal/logger"
)

type Service interface {
	WebhookProcessor
}

type service struct {
	WebhookProcessor
}

func NewService(logger logger.Logger, repo Repository, client Client) *service {
	return &service{
		WebhookProcessor: NewWebhookProcessor(logger, client, repo),
	}
}
