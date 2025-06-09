package upwardli

import (
	"context"
	"fmt"

	apiClient "elevate/packages/api-client-go"

	"github.com/pkg/errors"
)

// client types
type createWebhookRequest struct {
	WebhookName SubscriptionTopic `json:"webhook_name"`
	Endpoint    string            `json:"endpoint"`
}

type PartnerClient interface {
	GetAllWebhooks(ctx context.Context) ([]Webhook, error)
	CreateWebhook(ctx context.Context, webhook createWebhookRequest) (*Webhook, error)
	DeleteWebhook(ctx context.Context, webhookID string) error
}

type partnerClient struct {
	client   *apiClient.Client
	provider apiClient.TokenProvider
}

func NewPartnerClient(cfg partnerClientConfig) (PartnerClient, error) {
	if cfg.Scope == nil {
		defaultScope := "api:read api:write"
		cfg.Scope = &defaultScope
	}
	tokenProvider, err := NewPartnerTokenProvider(partnerClientConfig{
		Config: cfg.Config,
		Scope:  cfg.Scope,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize partner token provider")
	}

	newClient, err := apiClient.NewClient(cfg.BaseURL, apiClient.NewTokenAuthenticator(tokenProvider))
	if err != nil {
		return nil, err
	}
	return &partnerClient{
		client:   newClient,
		provider: tokenProvider,
	}, nil
}

func (c *partnerClient) GetAllWebhooks(ctx context.Context) ([]Webhook, error) {
	resp, err := c.client.Request(ctx, "/webhooks/registrations", apiClient.WithMethod(apiClient.MethodGet))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get webhooks")
	}

	var webhooksResp struct {
		Results []Webhook `external:"results"`
	}
	if err := UnmarshalExternal(resp, &webhooksResp); err != nil {
		return nil, errors.Wrap(err, "error parsing webhook response")
	}

	webhooks := make([]Webhook, len(webhooksResp.Results))
	copy(webhooks, webhooksResp.Results)

	return webhooks, nil
}

func (c *partnerClient) CreateWebhook(ctx context.Context, webhookReq createWebhookRequest) (*Webhook, error) {
	resp, err := c.client.Request(ctx, "/webhooks/registrations",
		apiClient.WithMethod(apiClient.MethodPost),
		apiClient.WithBody(webhookReq))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create webhook")
	}

	var respBody Webhook
	if err := UnmarshalExternal(resp, &respBody); err != nil {
		return nil, errors.Wrap(err, "error parsing webhook creation response")
	}
	return &respBody, nil
}

func (c *partnerClient) DeleteWebhook(ctx context.Context, webhookID string) error {
	if webhookID == "" {
		return errors.New("webhook ID is required")
	}

	_, err := c.client.Request(ctx, fmt.Sprintf("/webhooks/%s", webhookID),
		apiClient.WithMethod(apiClient.MethodDelete),
		apiClient.WithHeaders(map[string]string{"Content-Type": "application/json"}))
	if err != nil {
		return errors.Wrap(err, "failed to delete webhook")
	}

	return nil
}
