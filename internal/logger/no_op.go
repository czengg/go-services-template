package logger

import (
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

// NoOpLogger is a logger that does nothing
type NoOpLogger struct{}

func (n *NoOpLogger) Debug(msg string, fields ...zap.Field)             {}
func (n *NoOpLogger) Info(msg string, fields ...zap.Field)              {}
func (n *NoOpLogger) Warn(msg string, fields ...zap.Field)              {}
func (n *NoOpLogger) Error(msg string, fields ...zap.Field)             {}
func (n *NoOpLogger) Fatal(msg string, fields ...zap.Field)             {}
func (n *NoOpLogger) CaptureException(err error)                        {}
func (n *NoOpLogger) CaptureMessage(message string, level sentry.Level) {}
func (n *NoOpLogger) AddBreadcrumb(breadcrumb *sentry.Breadcrumb)       {}
func (n *NoOpLogger) With(fields ...zap.Field) Logger                   { return n }
func (n *NoOpLogger) WithUserID(userID string) Logger                   { return n }
func (n *NoOpLogger) WithRequestID(requestID string) Logger             { return n }
func (n *NoOpLogger) Close() error                                      { return nil }
