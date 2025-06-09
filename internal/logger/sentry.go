package logger

import (
	"time"

	"github.com/getsentry/sentry-go"
)

// initializeSentry sets up Sentry with the given configuration
func initializeSentry(config SentryConfig) error {
	return sentry.Init(sentry.ClientOptions{
		Dsn:         config.DSN,
		Environment: config.Environment,
		Release:     config.Release,
		Debug:       config.Debug,

		// Configure sampling
		SampleRate:       1.0,
		TracesSampleRate: 1.0,

		// Configure integrations
		Integrations: func(integrations []sentry.Integration) []sentry.Integration {
			// Remove default HTTP integration if not needed
			filteredIntegrations := []sentry.Integration{}
			for _, integration := range integrations {
				if integration.Name() != "Http" {
					filteredIntegrations = append(filteredIntegrations, integration)
				}
			}
			return filteredIntegrations
		},

		// Configure before send hook for filtering
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			// Filter out sensitive information
			if event.Request != nil {
				// Remove sensitive headers
				if event.Request.Headers != nil {
					delete(event.Request.Headers, "Authorization")
					delete(event.Request.Headers, "Cookie")
				}
			}
			return event
		},

		// Configure error filtering
		IgnoreErrors: []string{
			"context canceled",
			"context deadline exceeded",
		},
	})
}

// closeSentry flushes pending events and closes Sentry
func closeSentry(timeout time.Duration) {
	sentry.Flush(timeout)
}

// captureError captures an error to Sentry with additional context
func captureError(err error, tags map[string]string, extra map[string]interface{}) {
	sentry.WithScope(func(scope *sentry.Scope) {
		// Add tags
		for key, value := range tags {
			scope.SetTag(key, value)
		}

		// Add extra context
		for key, value := range extra {
			scope.SetExtra(key, value)
		}

		sentry.CaptureException(err)
	})
}

// captureMessage captures a message to Sentry with additional context
func captureMessage(message string, level sentry.Level, tags map[string]string, extra map[string]interface{}) {
	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(level)

		// Add tags
		for key, value := range tags {
			scope.SetTag(key, value)
		}

		// Add extra context
		for key, value := range extra {
			scope.SetExtra(key, value)
		}

		sentry.CaptureMessage(message)
	})
}

// addBreadcrumb adds a breadcrumb to the current scope
func addBreadcrumb(message, category string, level sentry.Level, data map[string]interface{}) {
	sentry.AddBreadcrumb(&sentry.Breadcrumb{
		Message:   message,
		Category:  category,
		Level:     level,
		Data:      data,
		Timestamp: time.Now(),
	})
}

// setUser sets the current user context for Sentry
func setUser(userID, role string) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{
			ID:   userID,
			Data: map[string]string{"role": role},
		})
	})
}

// setTag sets a tag in the current scope
func setTag(key, value string) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag(key, value)
	})
}

// setContext sets additional context in the current scope
func setContext(key string, context map[string]interface{}) {
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext(key, context)
	})
}
