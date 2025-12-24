package clientapi

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Anti-Pattern-Inc/listing-support-team-sample-sdk/generated/clientapi"
	"github.com/Anti-Pattern-Inc/listing-support-team-sample-sdk/middleware"
)

// getServerURL returns the API server URL from environment variable
func getServerURL() string {
	return os.Getenv("API_URL")
}

func withRequestEditorFns(ctx context.Context, c *clientapi.Client) error {
	// Get clientID and apiKey from environment variables
	clientID := os.Getenv("CLIENT_ID")
	apiKey := os.Getenv("API_KEY")

	if clientID == "" || apiKey == "" {
		return fmt.Errorf("environment variables CLIENT_ID and API_KEY are required")
	}

	// Authenticate to get token
	token, err := middleware.Authenticate(ctx, clientID, apiKey)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	c.RequestEditors = []clientapi.RequestEditorFn{
		func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+token)
			return nil
		},
	}

	return nil
}

// ClientWithResponse returns a ClientWithResponses with RequestEditorFn that generates authentication.
func ClientWithResponse(ctx context.Context) (*clientapi.ClientWithResponses, error) {
	serverURL := getServerURL()
	if serverURL == "" {
		return nil, fmt.Errorf("API_URL environment variable is required")
	}

	clientWithResponse, err := clientapi.NewClientWithResponses(serverURL, func(c *clientapi.Client) error {
		return withRequestEditorFns(ctx, c)
	})
	if err != nil {
		return nil, err
	}

	return clientWithResponse, nil
}
