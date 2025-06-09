package app

import (
	"template/internal/config"
	"template/internal/core/upwardli"
	"template/internal/logger"
)

type services struct {
	upwardli upwardli.Service
}

func newServices(config config.Config, logger logger.Logger, repos repositories, clients clients) services {
	upwardliService := upwardli.NewService(config.Upwardli(), logger, repos.Repository, clients.UpwardliPartner)
	if upwardliService == nil {
		logger.Fatal("failed to create upwardli service")
	}

	return services{
		upwardli: *upwardliService,
	}
}
