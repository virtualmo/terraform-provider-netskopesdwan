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

type gatewayListResponseDTO struct {
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

// ListGateways calls GET /v2/gateways and extracts gateway items without exposing
// response envelope assumptions outside the client.
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

	listResponse, err := decodeGatewayListResponse(body)
	if err != nil {
		return nil, fmt.Errorf("decode gateways response: %w", err)
	}

	return &GatewayListResult{
		Items: listResponse.Items,
	}, nil
}

func decodeGatewayListResponse(body []byte) (*gatewayListResponseDTO, error) {
	// TODO: Replace this temporary decoder with one confirmed /v2/gateways response
	// envelope once a real API sample is available. Keep all envelope uncertainty in
	// this client file so Terraform-facing code does not depend on guessed shapes.
	var arrayItems []json.RawMessage
	if err := json.Unmarshal(body, &arrayItems); err == nil {
		return &gatewayListResponseDTO{Items: arrayItems}, nil
	}

	var envelope struct {
		Data  []json.RawMessage `json:"data"`
		Items []json.RawMessage `json:"items"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, err
	}

	if len(envelope.Data) > 0 {
		return &gatewayListResponseDTO{Items: envelope.Data}, nil
	}

	if len(envelope.Items) > 0 {
		return &gatewayListResponseDTO{Items: envelope.Items}, nil
	}

	// TODO: Confirm the real empty-list envelope and whether pagination metadata is
	// returned alongside /v2/gateways results in the sample response.
	return &gatewayListResponseDTO{Items: []json.RawMessage{}}, nil
}
