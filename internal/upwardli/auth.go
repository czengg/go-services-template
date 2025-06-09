package upwardli

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	apiClient "elevate/packages/api-client-go"

	"github.com/pkg/errors"
)

type partnerClientConfig struct {
	Config
	Scope *string
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

func NewPartnerTokenProvider(config partnerClientConfig) (*partnerTokenProvider, error) {
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

// GetToken retrieves a partner authentication token, with caching to avoid unnecessary token requests.
func (p *partnerTokenProvider) GetToken(ctx context.Context) (string, error) {
	// Check if we have a valid cached token
	p.tokenMutex.RLock()
	if p.token != "" && time.Now().Before(p.expiresAt) {
		token := p.token
		p.tokenMutex.RUnlock()
		return token, nil
	}
	p.tokenMutex.RUnlock()

	// No valid token found, request a new one
	p.tokenMutex.Lock()
	defer p.tokenMutex.Unlock()

	// Double-check if token was refreshed by another goroutine while we were waiting
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

	// Store token with expiration time (subtracting a small buffer to ensure freshness)
	p.token = tokenResponse.AccessToken
	p.expiresAt = time.Now().Add(time.Duration(tokenResponse.ExpiresIn-60) * time.Second)

	return p.token, nil
}
