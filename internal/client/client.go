package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// GatewayListResult keeps the uncertain API response shape isolated inside the client.
type GatewayListResult struct {
	Items []json.RawMessage
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

// ListGateways calls GET /v2/gateways and extracts any gateway items without exposing
// assumptions about the response envelope outside the client.
func (c *Client) ListGateways(ctx context.Context) (*GatewayListResult, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, "/v2/gateways")
	if err != nil {
		return nil, fmt.Errorf("build gateways request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list gateways: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 4096))
		if readErr != nil {
			return nil, fmt.Errorf("list gateways: unexpected status %d and failed to read error body: %w", resp.StatusCode, readErr)
		}

		return nil, fmt.Errorf("list gateways: unexpected status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read gateways response: %w", err)
	}

	items, err := extractGatewayItems(body)
	if err != nil {
		return nil, fmt.Errorf("decode gateways response: %w", err)
	}

	return &GatewayListResult{
		Items: items,
	}, nil
}

func extractGatewayItems(body []byte) ([]json.RawMessage, error) {
	// TODO: Confirm the real /v2/gateways response envelope. This helper currently
	// accepts either a top-level JSON array or an object with a "data" or "items"
	// array so the uncertainty stays inside the client.
	var arrayItems []json.RawMessage
	if err := json.Unmarshal(body, &arrayItems); err == nil {
		return arrayItems, nil
	}

	var envelope struct {
		Data  []json.RawMessage `json:"data"`
		Items []json.RawMessage `json:"items"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}

	if len(envelope.Data) > 0 {
		return envelope.Data, nil
	}

	if len(envelope.Items) > 0 {
		return envelope.Items, nil
	}

	// TODO: Confirm whether an empty array can also be wrapped in another envelope
	// field name, and whether pagination metadata is returned alongside the list.
	return []json.RawMessage{}, nil
}
