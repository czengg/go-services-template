package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	DefaultTimeout = 30 * time.Second

	MethodGet    = http.MethodGet
	MethodPost   = http.MethodPost
	MethodPut    = http.MethodPut
	MethodDelete = http.MethodDelete
	MethodPatch  = http.MethodPatch
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	auth       Authenticator
}

type ClientOption func(*Client)

func NewClient(baseURL string, auth Authenticator, opts ...ClientOption) (*Client, error) {
	if baseURL == "" {
		return nil, ErrInvalidBaseURL
	}

	c := &Client{
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		baseURL: baseURL,
		auth:    auth,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func (c *Client) Request(ctx context.Context, path string, opts ...RequestOption) ([]byte, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	reqOpts := defaultRequestOptions()
	for _, opt := range opts {
		opt(reqOpts)
	}
	var reqURL = c.baseURL
	if reqOpts.subURL != "" {
		reqURL = reqOpts.subURL
	}
	fullURL, err := url.JoinPath(reqURL, path)
	if err != nil {
		return nil, errors.Wrap(err, "joining URL paths")
	}

	var bodyReader io.Reader
	if reqOpts.rawBody != nil {
		bodyReader = reqOpts.rawBody
	}
	if reqOpts.body != nil {
		jsonBody, err := json.Marshal(reqOpts.body)
		if err != nil {
			return nil, errors.Wrap(err, "marshaling request body")
		}
		bodyReader = bytes.NewBuffer(jsonBody)
		reqOpts.contentType = "application/json"
	}

	req, err := http.NewRequestWithContext(ctx, reqOpts.method, fullURL, bodyReader)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}

	req.Header.Set("Content-Type", reqOpts.contentType)
	for key, value := range reqOpts.headers {
		req.Header.Set(key, value)
	}

	if reqOpts.queryParams != nil {
		req.URL.RawQuery = reqOpts.queryParams.Encode()
	}

	if c.auth != nil {
		if err := c.auth.Authenticate(req); err != nil {
			return nil, errors.Wrap(err, "authenticating request")
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "executing request")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading response body")
	}

	if resp.StatusCode >= 400 {
		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			Body:       respBody,
		}
	}

	return respBody, nil
}

type RequestOptions struct {
	subURL      string
	method      string
	headers     map[string]string
	body        interface{}
	rawBody     io.Reader
	contentType string
	queryParams url.Values
	logger      *zap.Logger
}

func defaultRequestOptions() *RequestOptions {
	return &RequestOptions{
		method:  MethodGet,
		headers: make(map[string]string),
	}
}

type RequestOption func(*RequestOptions)

// WithSubURL sets an alternative base URL for a specific request.
// This overrides the client's default baseURL for this request only.
func WithSubURL(subURL string) RequestOption {
	return func(opts *RequestOptions) {
		opts.subURL = subURL
	}
}

// WithMethod sets the HTTP method for the request.
func WithMethod(method string) RequestOption {
	return func(opts *RequestOptions) {
		opts.method = method
	}
}

// WithHeaders adds headers to the request.
// It merges with existing headers, with new values overwriting existing ones.
func WithHeaders(headers map[string]string) RequestOption {
	return func(opts *RequestOptions) {
		for k, v := range headers {
			opts.headers[k] = v
		}
	}
}

// WithBody sets the request body.
// The body will be JSON-encoded before sending.
func WithBody(body interface{}) RequestOption {
	return func(opts *RequestOptions) {
		opts.body = body
	}
}

func WithBodyReader(r io.Reader, contentType string) RequestOption {
	return func(opts *RequestOptions) {
		opts.rawBody = r
		opts.contentType = contentType
	}
}

// WithQueryParams sets the query parameters.
// It replaces any existing query parameters.
func WithQueryParams(params url.Values) RequestOption {
	return func(opts *RequestOptions) {
		opts.queryParams = params
	}
}

// WithLogger sets the logger for the request.
func WithLogger(logger *zap.Logger) RequestOption {
	return func(opts *RequestOptions) {
		opts.logger = logger
	}
}

// HTTPError represents an error response from the server.
// It includes both the HTTP status code and response body.
type HTTPError struct {
	StatusCode int
	Body       []byte
}

// Error implements the error interface for HTTPError.
func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, string(e.Body))
}

// WithTimeout sets the client timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}
