package upwardli

import (
	"template/internal/logger"
)

type Service interface {
	WebhookProcessor
}

type service struct {
	WebhookProcessor
}

func NewService(config Config, logger logger.Logger, repo Repository, partnerClient PartnerClient) *service {
	return &service{
		WebhookProcessor: NewWebhookProcessor(logger, partnerClient, repo),
	}
}
