package clientapi

import (
	"context"
	"fmt"

	"github.com/Anti-Pattern-Inc/listing-support-team-sample-sdk/generated/clientapi"
)

// AuthClient provides a simple interface for authentication operations
type AuthClient struct {
	client *clientapi.ClientWithResponses
}

// NewAuthClient creates a new AuthClient instance
func NewAuthClient() (*AuthClient, error) {
	client, err := clientapi.NewClientWithResponses(server)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth client: %w", err)
	}

	return &AuthClient{
		client: client,
	}, nil
}

// GetAuthToken authenticates with clientID and apiKey, returns the access token
func (a *AuthClient) GetAuthToken(ctx context.Context, clientID, apiKey string) (string, error) {
	authRequest := clientapi.IssueAuthTokenJSONRequestBody{
		ClientId: clientID,
		ApiKey:   apiKey,
	}

	response, err := a.client.IssueAuthTokenWithResponse(ctx, authRequest)
	if err != nil {
		return "", fmt.Errorf("authentication request failed: %w", err)
	}

	switch response.StatusCode() {
	case 200:
		if response.JSON200 != nil && response.JSON200.AccessToken != nil {
			return *response.JSON200.AccessToken, nil
		}
		return "", fmt.Errorf("access token not found in response")
	case 401:
		return "", fmt.Errorf("invalid credentials: %s", response.JSON401.Error)
	case 500:
		return "", fmt.Errorf("server error: %s", response.JSON500.Error)
	default:
		return "", fmt.Errorf("unexpected error: status %d", response.StatusCode())
	}
}