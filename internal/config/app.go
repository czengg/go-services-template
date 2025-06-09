package config

import (
	"elevate/internal/db"
	"elevate/internal/logger"
	"elevate/internal/repository"
	"elevate/internal/router"
	"elevate/internal/service"

	"go.uber.org/zap"
)

type App struct {
	Server *router.Router
}

func NewApp(cfg Config, logger logger.Logger) *App {

	database, err := db.NewConnection(cfg.DB(), logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	repository := repository.NewRepository(database, logger)

	service := service.NewService(cfg, logger, repository)
	if service == nil {
		logger.Fatal("Failed to create service")
	}

	return &App{
		Server: router.NewRouter(*service),
	}
}
