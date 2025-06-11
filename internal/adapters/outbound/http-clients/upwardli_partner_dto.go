package httpclients

import (
	"template/internal/core/upwardli"
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

func (dto UpwardliWebhookDTO) ToDomain() upwardli.Webhook {
	return upwardli.Webhook{
		ID:          dto.ID,
		WebhookName: upwardli.SubscriptionTopic(dto.WebhookName),
		Endpoint:    dto.Endpoint,
		PartnerID:   dto.PartnerID,
		Status:      dto.Status,
		Failures:    dto.Failures,
		LastFailure: dto.LastFailure,
	}
}
