package app

import (
	httpclients "template/internal/adapters/outbound/http-clients"
	"template/internal/config"
	upwardli "template/internal/core/webhooks"
	"template/internal/logger"

	"go.uber.org/zap"
)

type clients struct {
	UpwardliPartner upwardli.Client
}

func newClients(config config.Config, logger logger.Logger) clients {
	upwardliPartner, err := httpclients.NewUpwardliPartnerClient(httpclients.UpwardliPartnerClientConfig{
		Config: config.Upwardli(),
		Scope:  nil,
	})
	if err != nil {
		logger.Fatal("failed to create upwardli partner client", zap.Error(err))
	}

	return clients{
		UpwardliPartner: upwardliPartner,
	}
}
