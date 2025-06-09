package app

import (
	"template/internal/adapters/outbound/http"
	"template/internal/config"
	"template/internal/core/upwardli"
	"template/internal/logger"

	"go.uber.org/zap"
)

type clients struct {
	UpwardliPartner upwardli.PartnerClient
}

func newClients(config config.Config, logger logger.Logger) *clients {
	upwardliPartner, err := http.NewUpwardliPartnerClient(http.UpwardliPartnerClientConfig{
		Config: config.Upwardli(),
		Scope:  nil,
	})
	if err != nil {
		logger.Fatal("failed to create upwardli partner client", zap.Error(err))
	}

	return &clients{
		UpwardliPartner: upwardliPartner,
	}
}
