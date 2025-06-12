package webhookprocessors

import (
	"context"
	"encoding/json"
	"fmt"
	httpclients "template/internal/adapters/outbound/http-clients"
	webhooks "template/internal/core/webhooks"
	"template/internal/logger"
	"time"
)

const (
	SubscriptionTopicConsumerCreated                  webhooks.SubscriptionTopic = "Consumer.Created"
	SubscriptionTopicConsumerUpdated                  webhooks.SubscriptionTopic = "Consumer.Updated"
	SubscriptionTopicConsumerClosed                   webhooks.SubscriptionTopic = "Consumer.Closed"
	SubscriptionTopicConsumerKYCStarted               webhooks.SubscriptionTopic = "Consumer.KYC.Started"
	SubscriptionTopicConsumerKYCPending               webhooks.SubscriptionTopic = "Consumer.KYC.Pending"
	SubscriptionTopicConsumerKYCCompleted             webhooks.SubscriptionTopic = "Consumer.KYC.Completed"
	SubscriptionTopicConsumerKYCNeedsReview           webhooks.SubscriptionTopic = "Consumer.KYC.NeedsReview"
	SubscriptionTopicConsumerKYCApproved              webhooks.SubscriptionTopic = "Consumer.KYC.Approved"
	SubscriptionTopicConsumerKYCFailed                webhooks.SubscriptionTopic = "Consumer.KYC.Failed"
	SubscriptionTopicPaymentCardCreated               webhooks.SubscriptionTopic = "PaymentCard.Created"
	SubscriptionTopicPaymentCardUpdated               webhooks.SubscriptionTopic = "PaymentCard.Updated"
	SubscriptionTopicPaymentCardClosed                webhooks.SubscriptionTopic = "PaymentCard.Closed"
	SubscriptionTopicPaymentCardTransactionSettlement webhooks.SubscriptionTopic = "PaymentCard.Transaction.Settlement"
	SubscriptionTopicACHSent                          webhooks.SubscriptionTopic = "ACH.Sent"
	SubscriptionTopicACHReceived                      webhooks.SubscriptionTopic = "ACH.Received"
	SubscriptionTopicACHFailed                        webhooks.SubscriptionTopic = "ACH.Failed"
	SubscriptionTopicPaymentTransferCreated           webhooks.SubscriptionTopic = "Payment.Transfer.Created"
	SubscriptionTopicPaymentTransferCompleted         webhooks.SubscriptionTopic = "Payment.Transfer.Completed"
	SubscriptionTopicPaymentTransferFailed            webhooks.SubscriptionTopic = "Payment.Transfer.Failed"
)

type upwardliWebhookEventRequest struct {
	ID              string                     `json:"id"`
	CreatedAt       *time.Time                 `json:"created_at"`
	EventName       webhooks.SubscriptionTopic `json:"event_name"`
	PartnerID       string                     `json:"partner_id"`
	Resources       []string                   `json:"resources"`
	LastAttemptedAt *time.Time                 `json:"last_attempted_at"`
	ResourcePath    string                     `json:"-"`
}

type upwardliProcessor struct {
	client httpclients.UpwardliPartnerClient
}

func NewUpwardliProcessor(l logger.Logger, client httpclients.UpwardliPartnerClient) webhooks.Processor {
	return &upwardliProcessor{
		client: client,
	}
}

func (p *upwardliProcessor) Process(ctx context.Context, body []byte, headers map[string]string) error {
	var req upwardliWebhookEventRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return err
	}

	entityInfo, err := p.client.GetEntityInfo(ctx, req.ResourcePath)
	if err != nil {
		return err
	}

	switch req.EventName {
	case SubscriptionTopicConsumerCreated:
		var dto httpclients.UpwardliConsumerDTO
		if err := json.Unmarshal(entityInfo, &dto); err != nil {
			return err
		}

		consumer := dto.ToDomain()
		fmt.Println(consumer)
		return nil
	}
	return nil
}
