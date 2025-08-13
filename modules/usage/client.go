package usage

import (
	"context"
	"net/http"
	"time"

	"github.com/nishimurashinya/openapi-go-sdk/client"
	"github.com/nishimurashinya/openapi-go-sdk/generated/usageapi"
	"github.com/nishimurashinya/openapi-go-sdk/middleware"
)

// Client provides high-level AWS Marketplace SaaS Usage Records API operations
type Client struct {
	baseClient *client.BaseClient
	usageAPI   *usageapi.ClientWithResponses
}

// NewClient creates a new usage API client with Bearer authentication
func NewClient(baseURL, apiKey string) (*Client, error) {
	auth := &middleware.AuthConfig{
		Type:   middleware.AuthTypeBearer,
		APIKey: apiKey,
	}
	
	return NewClientWithAuth(baseURL, auth)
}

// NewClientWithSignature creates a new usage API client with HMAC-SHA256 signature authentication
func NewClientWithSignature(baseURL, apiKey, secretKey string) (*Client, error) {
	auth := &middleware.AuthConfig{
		Type:      middleware.AuthTypeSignature,
		APIKey:    apiKey,
		SecretKey: secretKey,
	}
	
	return NewClientWithAuth(baseURL, auth)
}

// NewClientWithCustomAuth creates a new usage API client with custom authentication
func NewClientWithCustomAuth(baseURL string, customAuth func(ctx context.Context, req *http.Request) error) (*Client, error) {
	auth := &middleware.AuthConfig{
		Type:       middleware.AuthTypeCustom,
		CustomAuth: customAuth,
	}
	
	return NewClientWithAuth(baseURL, auth)
}

// NewClientWithAuth creates a new usage API client with provided authentication config
func NewClientWithAuth(baseURL string, auth *middleware.AuthConfig) (*Client, error) {
	baseClient := client.NewBaseClient(baseURL, auth)
	
	usageAPI, err := baseClient.CreateUsageAPIClient()
	if err != nil {
		return nil, err
	}
	
	return &Client{
		baseClient: baseClient,
		usageAPI:   usageAPI,
	}, nil
}

// V1API provides access to v1 usage record operations
type V1API struct {
	client *Client
}

// V1 returns a v1 API interface for usage records
func (c *Client) V1() *V1API {
	return &V1API{client: c}
}

// CreateMeteringUsageRecords creates multiple usage records with simplified interface
func (v *V1API) CreateMeteringUsageRecords(ctx context.Context, records []MeteringRecord) (*usageapi.CreateUsageRecordsResponse, error) {
	// Convert simplified records to internal format
	usageRecords := make([]usageapi.UsageRecord, len(records))
	for i, record := range records {
		usageRecord, err := convertMeteringRecord(record)
		if err != nil {
			return nil, err
		}
		usageRecords[i] = usageRecord
	}
	
	return v.client.usageAPI.CreateUsageRecordsWithResponse(ctx, usageRecords)
}

// UpdateMeteringUsageRecords updates multiple usage records
func (v *V1API) UpdateMeteringUsageRecords(ctx context.Context, records []MeteringRecord) (*usageapi.UpdateUsageRecordsResponse, error) {
	// Convert simplified records to internal format
	usageRecords := make([]usageapi.UsageRecord, len(records))
	for i, record := range records {
		usageRecord, err := convertMeteringRecord(record)
		if err != nil {
			return nil, err
		}
		usageRecords[i] = usageRecord
	}
	
	return v.client.usageAPI.UpdateUsageRecordsWithResponse(ctx, usageRecords)
}

// CreateUsageRecords creates new usage records (low-level API)
func (c *Client) CreateUsageRecords(ctx context.Context, records []usageapi.UsageRecord) (*usageapi.CreateUsageRecordsResponse, error) {
	return c.usageAPI.CreateUsageRecordsWithResponse(ctx, records)
}

// UpdateUsageRecords updates existing usage records (low-level API)
func (c *Client) UpdateUsageRecords(ctx context.Context, records []usageapi.UsageRecord) (*usageapi.UpdateUsageRecordsResponse, error) {
	return c.usageAPI.UpdateUsageRecordsWithResponse(ctx, records)
}

// OptionsUsageRecords performs CORS preflight request
func (c *Client) OptionsUsageRecords(ctx context.Context) (*usageapi.OptionsUsageRecordsResponse, error) {
	return c.usageAPI.OptionsUsageRecordsWithResponse(ctx)
}

// convertMeteringRecord converts a simplified MeteringRecord to internal UsageRecord
func convertMeteringRecord(record MeteringRecord) (usageapi.UsageRecord, error) {
	// Parse start time
	startTime, err := time.Parse(time.RFC3339, record.StartTime)
	if err != nil {
		return usageapi.UsageRecord{}, err
	}
	
	return usageapi.UsageRecord{
		ProductId:          record.ProductID,
		CustomerIdentifier: record.CustomerID,
		Dimension: usageapi.Dimension{
			Name:     record.Dimension.Name,
			Quantity: record.Dimension.Quantity,
		},
		StartTime: startTime,
	}, nil
}