package config

import (
	"elevate/internal/aws"
	"elevate/internal/db"
	"elevate/internal/plaid"
	"elevate/internal/upwardli"
	"os"
	"strings"

	"github.com/joho/godotenv"
	plaidSDK "github.com/plaid/plaid-go/v32/plaid"
)

func loadEnvFile() error {
	return godotenv.Load("config.env")
}

func loadFromEnv() (*config, error) {
	env := os.Getenv("ENV")
	if env == "" {
		env = "DEVELOPMENT"
	}

	local := false
	if strings.HasSuffix(env, "-LOCAL") {
		env = strings.Split(env, "-")[0]
		local = true
	}

	return &config{
		env:                  env,
		local:                local,
		sentryDSN:            os.Getenv("SENTRY_DSN"),
		interServiceSecret:   os.Getenv("INTER_SERVICE_SECRET"),
		clientJWTTokenSecret: os.Getenv("CLIENT_JWT_TOKEN_SECRET"),

		dbConfig: db.Config{
			User:     os.Getenv("DB_USER"),
			Endpoint: os.Getenv("DB_ENDPOINT"),
			Port:     os.Getenv("DB_PORT"),
			Password: os.Getenv("DB_PASSWORD"),
		},

		awsConfig: aws.Config{
			AccessID:     os.Getenv("AWS_ACCESS_KEY_ID"),
			AccessSecret: os.Getenv("AWS_SECRET_ACCESS_KEY"),
			Region:       os.Getenv("AWS_REGION"),
		},

		plaidConfig: plaid.Config{
			ClientID:      os.Getenv("PLAID_CLIENT_ID"),
			Secret:        os.Getenv("PLAID_SECRET"),
			Env:           plaidSDK.Environment(os.Getenv("PLAID_ENV")),
			SandboxSecret: os.Getenv("PLAID_SANDBOX_SECRET"),
			RedirectURL:   os.Getenv("PLAID_REDIRECT_URL"),
		},

		upwardliConfig: upwardli.Config{
			AuthURL:              os.Getenv("UPWARDLI_AUTH_URL"),
			APIURL:               os.Getenv("UPWARDLI_API_URL"),
			ClientID:             os.Getenv("UPWARDLI_CLIENT_ID"),
			ClientSecret:         os.Getenv("UPWARDLI_CLIENT_SECRET"),
			EmbeddedComponentURL: os.Getenv("UPWARDLI_EMBEDDED_COMPONENT_URL"),
			FBOAccountNumber:     os.Getenv("UPWARDLI_FBO_ACCOUNT_NUMBER"),
		},
	}, nil
}
