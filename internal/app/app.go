package app

import (
	"template/internal/adapters/outbound/persistence/mysql"
	"template/internal/config"
	"template/internal/logger"

	"go.uber.org/zap"
)

type App struct {
	Server router
}

func NewApp(cfg config.Config, logger logger.Logger) *App {

	database, err := mysql.NewConnection(cfg.DB(), logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	cronjobs := newCronJobs(logger)
	cronjobs.setupCronJobs()

	repos := newRepositories(database, logger)

	clients := newClients(cfg, logger)

	services := newServices(cfg, logger, repos, clients)

	router := newRouter(cfg, services)

	return &App{
		Server: router,
	}
}
