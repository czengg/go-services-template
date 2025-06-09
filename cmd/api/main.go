package main

import (
	"fmt"
	"log"
	"template/internal/app"
	"template/internal/config"
	"template/internal/logger"

	"go.uber.org/zap"
)

type Ctx struct {
	config config.Config
	logger logger.Logger
}

func main() {
	fmt.Println("My cabb - oh, forget it!!")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	loggerInstance, err := logger.New(logger.Config{
		Environment: cfg.Env(),
		ServiceName: "cabbage",
		SentryDSN:   cfg.SentryDSN(),
		Level:       getLogLevel(cfg),
		Local:       cfg.IsLocal(),
	})
	if err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer loggerInstance.Close()

	var c Ctx
	c.config = cfg
	c.logger = loggerInstance

	c.logger.Info("Starting service",
		zap.String("service", "cabbage"),
		zap.String("environment", cfg.Env()),
		zap.Bool("local", cfg.IsLocal()),
	)

	app := app.NewApp(cfg, c.logger)
	app.Server.Serve(cfg.Port(), c.logger)
}

func getLogLevel(cfg config.Config) string {
	if cfg.IsLocal() {
		return "debug"
	}
	if cfg.Env() == "PRODUCTION" {
		return "info"
	}
	return "debug"
}
