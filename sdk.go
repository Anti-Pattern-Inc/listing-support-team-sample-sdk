package sdk

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Config holds the configuration for the AWS Marketplace SaaS Usage Records SDK
type Config struct {
	// BaseURL is the base URL for the API (e.g., https://api-id.execute-api.ap-northeast-1.amazonaws.com/Prod)
	BaseURL string
	// APIKey is the Bearer token for authentication
	APIKey string
	// HTTPClient allows you to provide a custom HTTP client
	HTTPClient *http.Client
}

// SDK wraps the generated client with additional configuration and convenience methods
type SDK struct {
	client *ClientWithResponses
	config *Config
}

// NewSDK creates a new AWS Marketplace SaaS Usage Records SDK instance
func NewSDK(config *Config) (*SDK, error) {
	if config.BaseURL == "" {
		return nil, fmt.Errorf("BaseURL is required")
	}
	if config.APIKey == "" {
		return nil, fmt.Errorf("APIKey is required")
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	// Create request editor to add authentication
	requestEditor := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.APIKey))
		req.Header.Set("Content-Type", "application/json")
		return nil
	}

	client, err := NewClientWithResponses(config.BaseURL, WithHTTPClient(httpClient), WithRequestEditorFn(requestEditor))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return &SDK{
		client: client,
		config: config,
	}, nil
}

// CreateUsageRecords creates new usage records
func (s *SDK) CreateUsageRecords(ctx context.Context, records []UsageRecord) (*CreateUsageRecordsResponse, error) {
	return s.client.CreateUsageRecordsWithResponse(ctx, records)
}

// UpdateUsageRecords updates existing usage records
func (s *SDK) UpdateUsageRecords(ctx context.Context, records []UsageRecord) (*UpdateUsageRecordsResponse, error) {
	return s.client.UpdateUsageRecordsWithResponse(ctx, records)
}

// OptionsUsageRecords performs CORS preflight request
func (s *SDK) OptionsUsageRecords(ctx context.Context) (*OptionsUsageRecordsResponse, error) {
	return s.client.OptionsUsageRecordsWithResponse(ctx)
}

// Helper functions for creating records

// NewUsageRecord creates a new UsageRecord with required fields
func NewUsageRecord(productID, customerID, dimensionName string, quantity float32, startTime time.Time) UsageRecord {
	return UsageRecord{
		ProductId:          productID,
		CustomerIdentifier: customerID,
		Dimension: Dimension{
			Name:     dimensionName,
			Quantity: quantity,
		},
		StartTime: startTime,
	}
}

// NewUsageRecordWithAllocations creates a new UsageRecord with usage allocations
func NewUsageRecordWithAllocations(productID, customerID, dimensionName string, quantity float32, startTime time.Time, allocations []UsageAllocation) UsageRecord {
	record := NewUsageRecord(productID, customerID, dimensionName, quantity, startTime)
	record.Dimension.UsageAllocations = &allocations
	return record
}

// NewUsageAllocation creates a new UsageAllocation
func NewUsageAllocation(quantity float32, tags []Tag) UsageAllocation {
	return UsageAllocation{
		AllocatedUsageQuantity: &quantity,
		Tags:                   &tags,
	}
}

// NewTag creates a new Tag
func NewTag(key, value string) Tag {
	return Tag{
		Key:   key,
		Value: value,
	}
}