package xenditv3

import (
	"io"
	"net/http"
	"time"
)

// Client configures HTTP access to Xendit Payments API v3.
type Client struct {
	BaseURL      string
	SecretAPIKey string
	APIVersion   string
	HTTPClient   *http.Client
}

// NewClient returns a client with sandbox defaults.
func NewClient(secretAPIKey string) Client {
	return Client{
		BaseURL:      DefaultBaseURL,
		SecretAPIKey: secretAPIKey,
		APIVersion:   DefaultAPIVersion,
		HTTPClient:   &http.Client{Timeout: 80 * time.Second},
	}
}

func (c Client) apiVersion() string {
	if c.APIVersion != "" {
		return c.APIVersion
	}
	return DefaultAPIVersion
}

func (c Client) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}
	return &http.Client{Timeout: 80 * time.Second}
}

func (c Client) newPaymentsRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("api-version", c.apiVersion())
	req.SetBasicAuth(c.SecretAPIKey, "")
	return req, nil
}

func (c Client) newCustomerRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.SecretAPIKey, "")
	return req, nil
}

// Gateway exposes Payments API v3 operations.
type Gateway struct {
	Client Client
}

// NewGateway builds a gateway from client config.
func NewGateway(client Client) *Gateway {
	return &Gateway{Client: client}
}
