package webhooks

import "time"

// Internal types
type subscriptionTopic string

// whitelisted topics
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
	ID          string
	WebhookName subscriptionTopic
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Endpoint    string
	PartnerID   string
	Status      string
	Failures    int64
	LastFailure *time.Time
	Deleted     bool

	// only available when registering a webhook
	RegistrationID string
}
