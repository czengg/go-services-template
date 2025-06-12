package httpclients

import (
	banking "template/internal/core/banking"
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

		Provider: webhooks.ProviderUpwardli,
	}
}

type UpwardliConsumerDTO struct {
	ID            string   `json:"id,omitempty"` // Also referred to as ConsumerID (consumer_id) in some contexts
	PCID          string   `json:"pcid,omitempty"`
	ExternalID    string   `json:"external_id,omitempty"` // Also referred to as ConsumerID (consumer_id) in some contexts
	FirstName     string   `json:"first_name,omitempty"`
	LastName      string   `json:"last_name,omitempty"`
	Email         string   `json:"email,omitempty"`
	IsActive      bool     `json:"is_active,omitempty"`
	KYCStatus     string   `json:"kyc_status,omitempty"`
	PhoneNumber   string   `json:"phone_number,omitempty"`
	DateOfBirth   string   `json:"date_of_birth,omitempty"`
	TaxIDType     string   `json:"tax_id_type,omitempty"`
	TaxIdentifier string   `json:"tax_identifier,omitempty"`
	AddressLine1  string   `json:"address_line1,omitempty"`
	AddressLine2  string   `json:"address_line2,omitempty"`
	AddressCity   string   `json:"address_city,omitempty"`
	AddressState  string   `json:"address_state,omitempty"`
	AddressZip    string   `json:"address_zip,omitempty"`
	CreditLines   []string `json:"credit_lines,omitempty"`
}

func (dto UpwardliConsumerDTO) ToDomain() banking.Consumer {
	return banking.Consumer{
		ID:            dto.ID,
		PCID:          dto.PCID,
		ExternalID:    dto.ExternalID,
		IsActive:      dto.IsActive,
		KYCStatus:     dto.KYCStatus,
		TaxIDType:     dto.TaxIDType,
		TaxIdentifier: dto.TaxIdentifier,
	}
}
