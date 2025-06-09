package config

import (
	"fmt"
	"template/internal/adapters/outbound/persistence/mysql"
	"template/internal/core/aws"
	"template/internal/core/plaid"
	"template/internal/core/upwardli"
)

type Config interface {
	Env() string
	IsProduction() bool
	IsLocal() bool
	Port() string

	// Database configs
	DB() mysql.Config

	// External services
	AWS() aws.Config
	Plaid() plaid.Config
	Upwardli() upwardli.Config

	// Internal services
	InterServiceSecret() string
	ClientJWTTokenSecret() string

	// Other configs
	SentryDSN() string
}

type config struct {
	env                  string
	local                bool
	port                 string
	sentryDSN            string
	interServiceSecret   string
	clientJWTTokenSecret string
	dbConfig             mysql.Config
	awsConfig            aws.Config
	plaidConfig          plaid.Config
	upwardliConfig       upwardli.Config
}

func Load() (Config, error) {
	if err := loadEnvFile(); err != nil {
		fmt.Println("Failed to load .env file")
	}

	cfg, err := loadFromEnv()
	if err != nil {
		return nil, err
	}

	if err := validate(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *config) Env() string                  { return c.env }
func (c *config) IsProduction() bool           { return c.env == "PRODUCTION" }
func (c *config) IsLocal() bool                { return c.local }
func (c *config) Port() string                 { return c.port }
func (c *config) DB() mysql.Config             { return c.dbConfig }
func (c *config) AWS() aws.Config              { return c.awsConfig }
func (c *config) Plaid() plaid.Config          { return c.plaidConfig }
func (c *config) Upwardli() upwardli.Config    { return c.upwardliConfig }
func (c *config) InterServiceSecret() string   { return c.interServiceSecret }
func (c *config) ClientJWTTokenSecret() string { return c.clientJWTTokenSecret }
func (c *config) SentryDSN() string            { return c.sentryDSN }
