package app

import (
	webhookprocessors "template/internal/adapters/inbound/webhook-processors"
	webhooks "template/internal/core/webhooks"
	"template/internal/logger"
)

type webhookProcessors struct {
	UpwardliProcessor webhooks.Processor
}

func newWebhookProcessors(l logger.Logger, c clients) webhookProcessors {
	return webhookProcessors{
		UpwardliProcessor: webhookprocessors.NewUpwardliProcessor(l, c.UpwardliPartner),
	}
}
