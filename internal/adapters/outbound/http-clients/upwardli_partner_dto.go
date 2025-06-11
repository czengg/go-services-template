package httpclients

import (
	webhooks "template/internal/core/webhooks"
	"time"
)

type UpwardliWebhookDTO struct {
	ID          string     `json:"id"`
	WebhookName string     `json:"webhook_name"`
	Endpoint    string     `json:"endpoint"`
	PartnerID   string     `json:"partner_id"`
	Status      string     `json:"status"`
	Failures    int64      `json:"failures"`
	LastFailure *time.Time `json:"last_failure"`

	RegistrationID string `json:"registration_id,omitempty"`
}

func (dto UpwardliWebhookDTO) ToDomain() webhooks.Webhook {
	return webhooks.Webhook{
		ID:          dto.ID,
		WebhookName: webhooks.SubscriptionTopic(dto.WebhookName),
		Endpoint:    dto.Endpoint,
		PartnerID:   dto.PartnerID,
		Status:      dto.Status,
		Failures:    dto.Failures,
		LastFailure: dto.LastFailure,
	}
}
