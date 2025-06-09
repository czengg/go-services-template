package mysql

import (
	"fmt"
	"template/internal/logger"
	"time"

	"github.com/jmoiron/sqlx"
)

func NewConnection(cfg Config, logger logger.Logger) (*sqlx.DB, error) {
	logger.Info("Initializing MySQL connection...")

	connectionStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/", cfg.User, cfg.Password, cfg.Endpoint, cfg.Port)
	db, err := sqlx.Open("mysql", (connectionStr + "?parseTime=true"))
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 4)

	logger.Info("Opening MySQL connection...")
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	logger.Info("Successfully connected to MySQL.")
	return db, nil
}
