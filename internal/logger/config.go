package logger

// Config holds logger configuration
type Config struct {
	Environment string
	ServiceName string
	SentryDSN   string
	Level       string
	Local       bool
}

// SentryConfig holds Sentry-specific configuration
type SentryConfig struct {
	DSN         string
	Environment string
	Release     string
	Debug       bool
}
