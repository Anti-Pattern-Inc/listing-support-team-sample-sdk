package clientapi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Anti-Pattern-Inc/listing-support-team-sample-sdk/generated/clientapi"
)

var (
	server = "https://0zm201tv70.execute-api.ap-northeast-1.amazonaws.com/Prod"
)

// MarketplaceClient provides a simple interface for marketplace operations
type MarketplaceClient struct {
	Client *clientapi.ClientWithResponses
	token  string
}

// NewMarketplaceClient creates a new MarketplaceClient with Bearer token authentication
func NewMarketplaceClient(token string) (*MarketplaceClient, error) {
	authFunc := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	}

	client, err := clientapi.NewClientWithResponses(server, func(c *clientapi.Client) error {
		c.RequestEditors = []clientapi.RequestEditorFn{authFunc}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create marketplace client: %w", err)
	}

	return &MarketplaceClient{
		Client: client,
		token:  token,
	}, nil
}

// NewMarketplaceClientWithAuth creates a client by performing authentication first
func NewMarketplaceClientWithAuth(ctx context.Context, clientID, apiKey string) (*MarketplaceClient, error) {
	authClient, err := NewAuthClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create auth client: %w", err)
	}

	token, err := authClient.GetAuthToken(ctx, clientID, apiKey)
	if err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	return NewMarketplaceClient(token)
}