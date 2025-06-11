package repository

import (
	"template/internal/core/upwardli"
	"template/internal/logger"

	"template/internal/adapters/outbound/persistence/mysql/sqlc"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	upwardli.Repository
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
