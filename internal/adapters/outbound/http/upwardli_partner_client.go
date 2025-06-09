package http

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"template/internal/core/upwardli"
	apiClient "template/packages/api-client-go"
	"template/packages/common-go"

	"github.com/pkg/errors"
)

type UpwardliPartnerClientConfig struct {
	upwardli.Config
	Scope *string
}

type partnerClient struct {
	client   *apiClient.Client
	provider apiClient.TokenProvider
}

type partnerTokenProvider struct {
	BaseURL      string
	AuthURL      string
	GrantType    string
	ClientID     string
	ClientSecret string
	Scope        string
	Client       *apiClient.Client

	// Token caching fields
	token      string
	expiresAt  time.Time
	tokenMutex sync.RWMutex
}

type partnerTokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope"`
}

type partnerTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func NewUpwardliPartnerTokenProvider(config UpwardliPartnerClientConfig) (*partnerTokenProvider, error) {
	client, err := apiClient.NewClient(
		config.BaseURL,
		nil,
		apiClient.WithTimeout(30*time.Second),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize a new client with partner information")
	}

	if config.Scope == nil {
		defaultScope := "api:read api:write"
		config.Scope = &defaultScope
	}

	return &partnerTokenProvider{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		BaseURL:      config.BaseURL,
		AuthURL:      config.AuthURL,
		Client:       client,
		GrantType:    "client_credentials",
		Scope:        *config.Scope,
	}, nil
}

func (p *partnerTokenProvider) GetToken(ctx context.Context) (string, error) {
	p.tokenMutex.RLock()
	if p.token != "" && time.Now().Before(p.expiresAt) {
		token := p.token
		p.tokenMutex.RUnlock()
		return token, nil
	}
	p.tokenMutex.RUnlock()

	p.tokenMutex.Lock()
	defer p.tokenMutex.Unlock()

	if p.token != "" && time.Now().Before(p.expiresAt) {
		return p.token, nil
	}

	requestBody := partnerTokenRequest{
		GrantType:    p.GrantType,
		ClientID:     p.ClientID,
		ClientSecret: p.ClientSecret,
		Scope:        p.Scope,
	}

	respBody, err := p.Client.Request(ctx, "/auth/token/",
		apiClient.WithMethod(apiClient.MethodPost),
		apiClient.WithBody(requestBody),
		apiClient.WithSubURL(p.AuthURL),
	)
	if err != nil {
		return "", fmt.Errorf("requesting partner token: %w", err)
	}

	var tokenResponse partnerTokenResponse
	if err := json.Unmarshal(respBody, &tokenResponse); err != nil {
		return "", fmt.Errorf("parsing partner token response: %w", err)
	}

	p.token = tokenResponse.AccessToken
	p.expiresAt = time.Now().Add(time.Duration(tokenResponse.ExpiresIn-60) * time.Second)

	return p.token, nil
}

func NewUpwardliPartnerClient(cfg UpwardliPartnerClientConfig) (upwardli.PartnerClient, error) {
	if cfg.Scope == nil {
		defaultScope := "api:read api:write"
		cfg.Scope = &defaultScope
	}
	tokenProvider, err := NewUpwardliPartnerTokenProvider(UpwardliPartnerClientConfig{
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

func (c *partnerClient) GetAllWebhooks(ctx context.Context) ([]upwardli.Webhook, error) {
	resp, err := c.client.Request(ctx, "/webhooks/registrations", apiClient.WithMethod(apiClient.MethodGet))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get webhooks")
	}

	var webhooksResp struct {
		Results []upwardli.Webhook `external:"results"`
	}
	if err := common.UnmarshalExternal(resp, &webhooksResp); err != nil {
		return nil, errors.Wrap(err, "error parsing webhook response")
	}

	webhooks := make([]upwardli.Webhook, len(webhooksResp.Results))
	copy(webhooks, webhooksResp.Results)

	return webhooks, nil
}

func (c *partnerClient) CreateWebhook(ctx context.Context, webhookReq upwardli.CreateWebhookRequest) (*upwardli.Webhook, error) {
	resp, err := c.client.Request(ctx, "/webhooks/registrations",
		apiClient.WithMethod(apiClient.MethodPost),
		apiClient.WithBody(webhookReq))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create webhook")
	}

	var respBody upwardli.Webhook
	if err := common.UnmarshalExternal(resp, &respBody); err != nil {
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
