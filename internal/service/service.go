package service

import (
	"elevate/internal/config"
	"elevate/internal/logger"
	"elevate/internal/upwardli"
)

type Service struct {
	upwardli upwardli.Service
}

func NewService(config config.Config, logger logger.Logger, repo upwardli.Repository) *Service {
	upwardliService := upwardli.NewService(config.Upwardli(), logger, repo)
	if upwardliService == nil {
		logger.Fatal("failed to create upwardli service")
	}

	return &Service{
		upwardli: *upwardliService,
	}
}
