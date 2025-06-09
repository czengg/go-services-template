package cronjob

import "context"

type job interface {
	Execute(ctx context.Context) error
	GetID() string
	GetRetryCount() int
	IncrementRetry()
	GetMaxRetries() int
}
