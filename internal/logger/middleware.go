// internal/logger/middleware.go
package logger

import (
	"context"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type contextKey string

const (
	loggerContextKey    contextKey = "logger"
	requestIDContextKey contextKey = "request_id"
)

// Middleware returns an HTTP middleware that adds logging context
func Middleware(logger Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Generate request ID
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.New().String()
			}

			// Create request-scoped logger
			requestLogger := logger.WithRequestID(requestID).With(
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
				zap.String("user_agent", r.UserAgent()),
			).WithUserID(r.Context().Value("userID").(string))

			// Add to context
			ctx := context.WithValue(r.Context(), loggerContextKey, requestLogger)
			ctx = context.WithValue(ctx, requestIDContextKey, requestID)
			r = r.WithContext(ctx)

			// Add request ID to response headers
			w.Header().Set("X-Request-ID", requestID)

			// Wrap response writer to capture status code
			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Log request start
			requestLogger.Info("Request started")

			// Add breadcrumb for Sentry
			requestLogger.AddBreadcrumb(&sentry.Breadcrumb{
				Message:  "HTTP Request",
				Category: "http",
				Level:    sentry.LevelInfo,
				Data: map[string]interface{}{
					"method":     r.Method,
					"path":       r.URL.Path,
					"request_id": requestID,
				},
			})

			// Process request
			next.ServeHTTP(wrapped, r)

			// Calculate duration
			duration := time.Since(start)

			// Log request completion
			logLevel := getLogLevelFromStatus(wrapped.statusCode)
			fields := []zap.Field{
				zap.Int("status_code", wrapped.statusCode),
				zap.Duration("duration", duration),
				zap.Int64("duration_ms", duration.Milliseconds()),
			}

			switch logLevel {
			case "debug":
				requestLogger.Debug("Request completed", fields...)
			case "info":
				requestLogger.Info("Request completed", fields...)
			case "warn":
				requestLogger.Warn("Request completed", fields...)
			case "error":
				requestLogger.Error("Request completed", fields...)
			}
		})
	}
}

// fromContext extracts the logger from the request context
func fromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerContextKey).(Logger); ok {
		return logger
	}
	// Return a no-op logger if not found
	return &NoOpLogger{}
}

// requestIDFromContext extracts the request ID from the context
func requestIDFromContext(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDContextKey).(string); ok {
		return requestID
	}
	return ""
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriter) Write(data []byte) (int, error) {
	return w.ResponseWriter.Write(data)
}

// getLogLevelFromStatus returns appropriate log level based on HTTP status
func getLogLevelFromStatus(statusCode int) string {
	switch {
	case statusCode >= 500:
		return "error"
	case statusCode >= 400:
		return "warn"
	case statusCode >= 300:
		return "info"
	default:
		return "info"
	}
}
