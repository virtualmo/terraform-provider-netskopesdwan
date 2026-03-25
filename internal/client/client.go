package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Client is a minimal handwritten API client for Netskope SD-WAN API v2.
type Client struct {
	baseURL    string
	apiToken   string
	httpClient *http.Client
	insecure   bool
}

// New creates a client with a small, explicit surface area.
func New(baseURL, apiToken string, timeout time.Duration, insecure bool) *Client {
	return &Client{
		baseURL:  strings.TrimRight(baseURL, "/"),
		apiToken: apiToken,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		insecure: insecure,
	}
}

// NewRequest builds an authenticated HTTP request against the configured base URL.
func (c *Client) NewRequest(ctx context.Context, method, path string) (*http.Request, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiToken)
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// Do executes the request with the shared HTTP client.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// TODO: Honor insecure once the transport requirements are confirmed.
	_ = c.insecure

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	return resp, nil
}
