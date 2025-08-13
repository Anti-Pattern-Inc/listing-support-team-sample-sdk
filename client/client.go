package client

import (
	"fmt"
	"net/http"

	"github.com/nishimurashinya/openapi-go-sdk/generated/usageapi"
	"github.com/nishimurashinya/openapi-go-sdk/middleware"
)

// BaseClient provides common functionality for all AWS Marketplace SaaS Usage Records API clients
type BaseClient struct {
	BaseURL    string
	Auth       *middleware.AuthConfig
	HTTPClient *http.Client
}

// NewBaseClient creates a new base client with authentication
func NewBaseClient(baseURL string, auth *middleware.AuthConfig) *BaseClient {
	return &BaseClient{
		BaseURL: baseURL,
		Auth:    auth,
		HTTPClient: &http.Client{},
	}
}

// WithHTTPClient sets a custom HTTP client
func (c *BaseClient) WithHTTPClient(client *http.Client) *BaseClient {
	c.HTTPClient = client
	return c
}

// CreateUsageAPIClient creates a configured usage API client
func (c *BaseClient) CreateUsageAPIClient() (*usageapi.ClientWithResponses, error) {
	// Validate configuration
	if c.BaseURL == "" {
		return nil, fmt.Errorf("BaseURL is required")
	}
	
	if c.Auth == nil {
		return nil, fmt.Errorf("authentication configuration is required")
	}
	
	if err := middleware.ValidateAuthConfig(c.Auth); err != nil {
		return nil, fmt.Errorf("invalid auth config: %w", err)
	}
	
	// Create request editor for authentication
	requestEditor, err := c.Auth.CreateRequestEditor()
	if err != nil {
		return nil, fmt.Errorf("failed to create auth handler: %w", err)
	}
	
	// Create and return the usage API client
	client, err := usageapi.NewClientWithResponses(
		c.BaseURL,
		usageapi.WithHTTPClient(c.HTTPClient),
		usageapi.WithRequestEditorFn(requestEditor),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create usage API client: %w", err)
	}
	
	return client, nil
}