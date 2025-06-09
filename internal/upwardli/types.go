package upwardli

// Webhook types
type subscriptionTopic string

const (
	SubscriptionTopicConsumerCreated                  subscriptionTopic = "Consumer.Created"
	SubscriptionTopicConsumerUpdated                  subscriptionTopic = "Consumer.Updated"
	SubscriptionTopicConsumerClosed                   subscriptionTopic = "Consumer.Closed"
	SubscriptionTopicConsumerKYCStarted               subscriptionTopic = "Consumer.KYC.Started"
	SubscriptionTopicConsumerKYCPending               subscriptionTopic = "Consumer.KYC.Pending"
	SubscriptionTopicConsumerKYCCompleted             subscriptionTopic = "Consumer.KYC.Completed"
	SubscriptionTopicConsumerKYCNeedsReview           subscriptionTopic = "Consumer.KYC.NeedsReview"
	SubscriptionTopicConsumerKYCApproved              subscriptionTopic = "Consumer.KYC.Approved"
	SubscriptionTopicConsumerKYCFailed                subscriptionTopic = "Consumer.KYC.Failed"
	SubscriptionTopicPaymentCardCreated               subscriptionTopic = "PaymentCard.Created"
	SubscriptionTopicPaymentCardUpdated               subscriptionTopic = "PaymentCard.Updated"
	SubscriptionTopicPaymentCardClosed                subscriptionTopic = "PaymentCard.Closed"
	SubscriptionTopicPaymentCardTransactionSettlement subscriptionTopic = "PaymentCard.Transaction.Settlement"
	SubscriptionTopicACHSent                          subscriptionTopic = "ACH.Sent"
	SubscriptionTopicACHReceived                      subscriptionTopic = "ACH.Received"
	SubscriptionTopicACHFailed                        subscriptionTopic = "ACH.Failed"
	SubscriptionTopicPaymentTransferCreated           subscriptionTopic = "Payment.Transfer.Created"
	SubscriptionTopicPaymentTransferCompleted         subscriptionTopic = "Payment.Transfer.Completed"
	SubscriptionTopicPaymentTransferFailed            subscriptionTopic = "Payment.Transfer.Failed"
)

type webhook struct {
	ID          string            `json:"id" db:"id" external:"id"`
	WebhookName subscriptionTopic `json:"webhookName" db:"webhook_name" external:"webhook_name"`
	CreatedAt   string            `json:"createdAt" db:"created_at" external:"-"`
	UpdatedAt   string            `json:"updatedAt" db:"updated_at" external:"-"`
	Endpoint    string            `json:"endpoint" db:"endpoint" external:"endpoint"`
	PartnerID   string            `json:"partnerId" db:"partner_id" external:"partner_id"`
	Status      string            `json:"status" db:"status" external:"status"`
	Failures    int64             `json:"failures" db:"failures" external:"failures"`
	LastFailure string            `json:"lastFailure" db:"last_failure" external:"last_failure"`
	Deleted     bool              `json:"deleted" db:"deleted" external:"-"`

	// only available when registering a webhook
	RegistrationID string `json:"registrationId" db:"-" external:"registration_id"`
}
