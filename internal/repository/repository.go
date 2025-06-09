package repository

import (
	"elevate/internal/logger"
	"elevate/internal/upwardli"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	upwardli.Repository
}

type repository struct {
	db     *sqlx.DB
	logger logger.Logger
}

func NewRepository(db *sqlx.DB, logger logger.Logger) Repository {
	return &repository{
		db:     db,
		logger: logger,
	}
}
