package repository

import (
	webhooks "template/internal/core/webhooks"
	"template/internal/logger"

	"template/internal/adapters/outbound/persistence/mysql/sqlc"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	webhooks.Repository
}

type repository struct {
	db      *sqlx.DB
	logger  logger.Logger
	queries *sqlc.Queries
}

func NewRepository(db *sqlx.DB, logger logger.Logger) Repository {
	return &repository{
		db:      db,
		logger:  logger,
		queries: sqlc.New(db),
	}
}
