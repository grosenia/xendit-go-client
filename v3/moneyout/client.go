package moneyout

import (
	"io"
	"net/http"
	"time"
)

const DefaultBaseURL = "https://api.xendit.co"

// Client configures HTTP access to Xendit Money Out APIs (payout, batch disbursement).
// These are NOT Payments API v3 — same legacy disbursement endpoints.
type Client struct {
	BaseURL      string
	SecretAPIKey string
	HTTPClient   *http.Client
}

// NewClient returns a client with defaults.
func NewClient(secretAPIKey string) Client {
	return Client{
		BaseURL:      DefaultBaseURL,
		SecretAPIKey: secretAPIKey,
		HTTPClient:   &http.Client{Timeout: 80 * time.Second},
	}
}

func (c Client) baseURL() string {
	if c.BaseURL != "" {
		return c.BaseURL
	}
	return DefaultBaseURL
}

func (c Client) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return &http.Client{Timeout: 80 * time.Second}
}

func (c Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, c.baseURL()+path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.SecretAPIKey, "")
	return req, nil
}

func (c Client) newBatchRequest(idempotencyKey, method, path string, body io.Reader) (*http.Request, error) {
	req, err := c.newRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-IDEMPOTENCY-KEY", idempotencyKey)
	return req, nil
}

// Gateway exposes Money Out operations.
type Gateway struct {
	Client Client
}

// NewGateway builds a gateway from client config.
func NewGateway(client Client) *Gateway {
	return &Gateway{Client: client}
}
