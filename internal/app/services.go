package app

import (
	"template/internal/config"
	webhooks "template/internal/core/webhooks"
	"template/internal/logger"
)

type services struct {
	webhooks webhooks.Service
}

func newServices(config config.Config, logger logger.Logger, repos repositories, clients clients) services {
	webhooksService := webhooks.NewService(logger, repos.Repository, clients.UpwardliPartner, webhooks.ProviderUpwardli)
	if webhooksService == nil {
		logger.Fatal("failed to create upwardli service")
	}

	return services{
		webhooks: *webhooksService,
	}
}
