package middleware

import (
	"context"
	"fmt"
	"os"

	"github.com/Anti-Pattern-Inc/listing-support-team-sample-sdk/generated/clientapi"
)

// getServerURL returns the API server URL from environment variable
func getServerURL() string {
	return os.Getenv("API_URL")
}

// AuthClient provides authentication functionality
type AuthClient struct {
	client *clientapi.ClientWithResponses
}

// NewAuthClient creates a new AuthClient instance
func NewAuthClient() (*AuthClient, error) {
	serverURL := getServerURL()
	if serverURL == "" {
		return nil, fmt.Errorf("API_URL environment variable is required")
	}

	client, err := clientapi.NewClientWithResponses(serverURL)
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

	return a.handleAuthResponse(response)
}

// handleAuthResponse handles the authentication response and extracts the token
func (a *AuthClient) handleAuthResponse(response *clientapi.IssueAuthTokenResponse) (string, error) {
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

// Authenticate is a convenience function that creates an AuthClient and gets a token
// If clientID or apiKey is empty, it will try to get them from environment variables
func Authenticate(ctx context.Context, clientID, apiKey string) (string, error) {
	// If clientID or apiKey is empty, try to get from environment variables
	if clientID == "" {
		clientID = os.Getenv("CLIENT_ID")
	}
	if apiKey == "" {
		apiKey = os.Getenv("API_KEY")
	}

	if clientID == "" || apiKey == "" {
		return "", fmt.Errorf("clientID and apiKey are required (either as parameters or environment variables CLIENT_ID and API_KEY)")
	}

	authClient, err := NewAuthClient()
	if err != nil {
		return "", err
	}

	return authClient.GetAuthToken(ctx, clientID, apiKey)
}
