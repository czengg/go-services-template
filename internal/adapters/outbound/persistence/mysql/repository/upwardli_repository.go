package repository

import (
	"context"
	"template/internal/adapters/outbound/persistence/mysql/sqlc"
	banking "template/internal/core/banking"
)

func (r *repository) SaveBankingConsumer(ctx context.Context, consumer banking.Consumer) error {
	return r.queries.SaveUpwardliConsumer(ctx, sqlc.SaveUpwardliConsumerParams{
		ID:         consumer.ID,
		Pcid:       consumer.PCID,
		ExternalID: consumer.ExternalID,
		IsActive:   consumer.IsActive,
		KycStatus:  consumer.KYCStatus,
		TaxIDType:  consumer.TaxIDType,
	})
}
