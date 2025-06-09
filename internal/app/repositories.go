package app

import (
	"template/internal/adapters/outbound/persistence/repository"
	"template/internal/logger"

	"github.com/jmoiron/sqlx"
)

type repositories struct {
	Repository repository.Repository
}

func newRepositories(db *sqlx.DB, logger logger.Logger) repositories {
	repo := repository.NewRepository(db, logger)
	if repo == nil {
		logger.Fatal("failed to create repository")
	}

	return repositories{
		Repository: repo,
	}
}
