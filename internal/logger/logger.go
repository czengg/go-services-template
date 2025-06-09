package logger

import (
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger defines the interface for application logging
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)

	// Sentry-specific methods
	CaptureException(err error)
	CaptureMessage(message string, level sentry.Level)
	AddBreadcrumb(breadcrumb *sentry.Breadcrumb)

	// Context methods
	With(fields ...zap.Field) Logger
	WithUserID(userID string) Logger
	WithRequestID(requestID string) Logger

	// Close flushes any pending logs
	Close() error
}

type logger struct {
	zap        *zap.Logger
	config     Config
	baseFields []zap.Field
}

// New creates a new logger instance with Sentry integration
func New(config Config) (Logger, error) {
	// Initialize Sentry first
	if config.SentryDSN != "" {
		sentryConfig := SentryConfig{
			DSN:         config.SentryDSN,
			Environment: config.Environment,
			Debug:       config.Local,
		}

		if err := initializeSentry(sentryConfig); err != nil {
			return nil, fmt.Errorf("failed to initialize Sentry: %w", err)
		}
	}

	// Create zap logger
	zapLogger, err := createZapLogger(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create zap logger: %w", err)
	}

	// Set Sentry tags for the service
	if config.SentryDSN != "" {
		setTag("service", config.ServiceName)
		setTag("environment", config.Environment)
	}

	return &logger{
		zap:    zapLogger,
		config: config,
		baseFields: []zap.Field{
			zap.String("service", config.ServiceName),
			zap.String("environment", config.Environment),
		},
	}, nil
}

// createZapLogger creates and configures a zap logger
func createZapLogger(config Config) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	if config.Level != "" {
		if err := level.UnmarshalText([]byte(config.Level)); err != nil {
			return nil, fmt.Errorf("invalid log level %s: %w", config.Level, err)
		}
	}

	zapConfig := zap.Config{
		Level:       zap.NewAtomicLevelAt(level),
		Development: config.Local,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	if config.Local {
		// Pretty output for local development
		zapConfig.Encoding = "console"
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
	}

	return zapConfig.Build()
}

// Debug logs a debug message
func (l *logger) Debug(msg string, fields ...zap.Field) {
	allFields := append(l.baseFields, fields...)
	l.zap.Debug(msg, allFields...)
}

// Info logs an info message
func (l *logger) Info(msg string, fields ...zap.Field) {
	allFields := append(l.baseFields, fields...)
	l.zap.Info(msg, allFields...)
}

// Warn logs a warning message
func (l *logger) Warn(msg string, fields ...zap.Field) {
	allFields := append(l.baseFields, fields...)
	l.zap.Warn(msg, allFields...)

	// Add breadcrumb to Sentry
	if l.config.SentryDSN != "" {
		l.AddBreadcrumb(&sentry.Breadcrumb{
			Message:  msg,
			Level:    sentry.LevelWarning,
			Category: "log",
			Data:     fieldsToMap(fields),
		})
	}
}

// Error logs an error message and sends to Sentry
func (l *logger) Error(msg string, fields ...zap.Field) {
	allFields := append(l.baseFields, fields...)
	l.zap.Error(msg, allFields...)

	// Send to Sentry if configured
	if l.config.SentryDSN != "" {
		captureMessage(msg, sentry.LevelError, nil, fieldsToMap(fields))
	}
}

// Fatal logs a fatal message and sends to Sentry
func (l *logger) Fatal(msg string, fields ...zap.Field) {
	allFields := append(l.baseFields, fields...)
	l.zap.Fatal(msg, allFields...)

	// Send to Sentry if configured
	if l.config.SentryDSN != "" {
		captureMessage(msg, sentry.LevelFatal, nil, fieldsToMap(fields))
		closeSentry(2 * time.Second) // Wait for Sentry to send
	}
}

// CaptureException captures an exception to Sentry
func (l *logger) CaptureException(err error) {
	if l.config.SentryDSN != "" {
		sentry.CaptureException(err)
	}
	l.Error("Exception captured", zap.Error(err))
}

// CaptureMessage captures a message to Sentry
func (l *logger) CaptureMessage(message string, level sentry.Level) {
	if l.config.SentryDSN != "" {
		sentry.CaptureMessage(message)
	}

	// Also log locally
	switch level {
	case sentry.LevelDebug:
		l.Debug(message)
	case sentry.LevelInfo:
		l.Info(message)
	case sentry.LevelWarning:
		l.Warn(message)
	case sentry.LevelError:
		l.Error(message)
	case sentry.LevelFatal:
		l.Fatal(message)
	}
}

// AddBreadcrumb adds a breadcrumb to Sentry
func (l *logger) AddBreadcrumb(breadcrumb *sentry.Breadcrumb) {
	if l.config.SentryDSN != "" {
		sentry.AddBreadcrumb(breadcrumb)
	}
}

// With returns a new logger with additional fields
func (l *logger) With(fields ...zap.Field) Logger {
	newBaseFields := append(l.baseFields, fields...)
	return &logger{
		zap:        l.zap,
		config:     l.config,
		baseFields: newBaseFields,
	}
}

// WithUserID returns a new logger with user ID context
func (l *logger) WithUserID(userID string) Logger {
	// Set user in Sentry scope
	if l.config.SentryDSN != "" {
		setUser(userID, "", "")
	}

	return l.With(zap.String("user_id", userID))
}

// WithRequestID returns a new logger with request ID context
func (l *logger) WithRequestID(requestID string) Logger {
	return l.With(zap.String("request_id", requestID))
}

// Close flushes any pending logs and closes Sentry
func (l *logger) Close() error {
	if l.config.SentryDSN != "" {
		closeSentry(2 * time.Second)
	}
	return l.zap.Sync()
}

// Helper function to convert zap fields to map for Sentry
func fieldsToMap(fields []zap.Field) map[string]interface{} {
	result := make(map[string]interface{})
	for _, field := range fields {
		switch field.Type {
		case zapcore.StringType:
			result[field.Key] = field.String
		case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type:
			result[field.Key] = field.Integer
		case zapcore.Float64Type, zapcore.Float32Type:
			result[field.Key] = field.Integer // Note: zap stores floats as integers
		case zapcore.BoolType:
			result[field.Key] = field.Integer == 1
		default:
			result[field.Key] = field.Interface
		}
	}
	return result
}
