package client

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Authenticator interface {
	Authenticate(*http.Request) error
}

type TokenProvider interface {
	GetToken(context.Context) (string, error)
}

type TokenAuthenticator struct {
	provider TokenProvider
	token    string
	mu       sync.RWMutex
	expires  time.Time
}

func NewTokenAuthenticator(provider TokenProvider) *TokenAuthenticator {
	return &TokenAuthenticator{
		provider: provider,
	}
}

func (a *TokenAuthenticator) Authenticate(req *http.Request) error {
	token, err := a.getValidToken(req.Context())
	if err != nil {
		return fmt.Errorf("getting valid token: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return nil
}

func (a *TokenAuthenticator) getValidToken(ctx context.Context) (string, error) {
	a.mu.RLock()
	if a.isTokenValid() {
		token := a.token
		a.mu.RUnlock()
		return token, nil
	}
	a.mu.RUnlock()

	// Need to refresh token
	a.mu.Lock()
	defer a.mu.Unlock()

	// Double check after acquiring write lock
	if a.isTokenValid() {
		return a.token, nil
	}

	token, err := a.provider.GetToken(ctx)
	if err != nil {
		return "", fmt.Errorf("refreshing token: %w", err)
	}

	a.token = token
	a.expires = time.Now().Add(55 * time.Minute)
	return token, nil
}

func (a *TokenAuthenticator) isTokenValid() bool {
	return a.token != "" && time.Now().Before(a.expires)
}

type BasicAuthenticator struct {
	username string
	password string
}

func NewBasicAuthenticator(username, password string) *BasicAuthenticator {
	return &BasicAuthenticator{
		username: username,
		password: password,
	}
}

func (a *BasicAuthenticator) Authenticate(req *http.Request) error {
	req.SetBasicAuth(a.username, a.password)
	return nil
}
