package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

// GatewayListResult contains the confirmed /v2/gateways response fields used by the provider.
type GatewayListResult struct {
	PageInfo GatewayPageInfoDTO
	Items    []GatewayDTO
}

type gatewayListResponseDTO struct {
	PageInfo GatewayPageInfoDTO `json:"page_info"`
	Data     []GatewayDTO       `json:"data"`
}

type GatewayPageInfoDTO struct {
	EndCursor  string `json:"end_cursor"`
	HasNext    bool   `json:"has_next"`
	TotalCount int64  `json:"total_count"`
}

type GatewayDTO struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	CreatedAt            string `json:"created_at"`
	ModifiedAt           string `json:"modified_at"`
	ConfigUpdatesEnabled bool   `json:"config_updates_enabled"`
	Managed              bool   `json:"managed"`
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
		PageInfo: listResponse.PageInfo,
		Items:    listResponse.Data,
	}, nil
}

// GetGateway calls GET /v2/gateways/{id} and decodes the confirmed single-gateway response.
func (c *Client) GetGateway(ctx context.Context, id string) (*GatewayDTO, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, "/v2/gateways/"+url.PathEscape(id))
	if err != nil {
		return nil, fmt.Errorf("build gateway request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get gateway: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 4096))
		if readErr != nil {
			return nil, fmt.Errorf("get gateway: unexpected status %d and failed to read error body: %w", resp.StatusCode, readErr)
		}

		return nil, fmt.Errorf("get gateway: unexpected status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read gateway response: %w", err)
	}

	gateway, err := decodeGatewayResponse(body)
	if err != nil {
		return nil, fmt.Errorf("decode gateway response: %w", err)
	}

	return gateway, nil
}

func decodeGatewayListResponse(body []byte) (*gatewayListResponseDTO, error) {
	var response gatewayListResponseDTO
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	if response.Data == nil {
		response.Data = []GatewayDTO{}
	}

	return &response, nil
}

func decodeGatewayResponse(body []byte) (*GatewayDTO, error) {
	var response GatewayDTO
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
