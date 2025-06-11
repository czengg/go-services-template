package httphandlers

import (
	"template/internal/core/upwardli"
	"time"
)

type WebhookResponse struct {
	ID             string  `json:"id"`
	WebhookName    string  `json:"webhookName"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
	Endpoint       string  `json:"endpoint"`
	PartnerID      string  `json:"partnerId"`
	Status         string  `json:"status"`
	Failures       int64   `json:"failures"`
	LastFailure    *string `json:"lastFailure,omitempty"`
	RegistrationID string  `json:"registrationId,omitempty"`
}

func WebhookToResponse(w upwardli.Webhook) WebhookResponse {
	resp := WebhookResponse{
		ID:          w.ID,
		WebhookName: string(w.WebhookName),
		CreatedAt:   w.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   w.UpdatedAt.Format(time.RFC3339),
		Endpoint:    w.Endpoint,
		PartnerID:   w.PartnerID,
		Status:      w.Status,
		Failures:    w.Failures,
	}

	if w.LastFailure != nil {
		lastFailure := w.LastFailure.Format(time.RFC3339)
		resp.LastFailure = &lastFailure
	}

	return resp
}
