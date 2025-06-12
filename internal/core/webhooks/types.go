package webhooks

import "time"

// Internal types
type subscriptionTopic string

type provider string

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

	Provider provider

	// only available when registering a webhook
	RegistrationID string
}
