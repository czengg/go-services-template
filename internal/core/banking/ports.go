package banking

import "context"

type ConsumerManager interface {
	SaveConsumer(ctx context.Context, consumer Consumer) error
}

type Repository interface {
	SaveBankingConsumer(ctx context.Context, consumer Consumer) error
}

type Consumer = consumer
