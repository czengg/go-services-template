package banking

import "context"

type consumerManager struct {
	repo Repository
}

func NewConsumerManager(repo Repository) ConsumerManager {
	return &consumerManager{
		repo: repo,
	}
}

func (m *consumerManager) SaveConsumer(ctx context.Context, consumer Consumer) error {
	return m.repo.SaveBankingConsumer(ctx, consumer)
}
