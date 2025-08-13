package usage

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/nishimurashinya/openapi-go-sdk/generated/usageapi"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient("https://api.example.com", "test-api-key")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	
	if client == nil {
		t.Fatal("Client should not be nil")
	}
}

func TestNewClientWithSignature(t *testing.T) {
	client, err := NewClientWithSignature("https://api.example.com", "test-api-key", "test-secret")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	
	if client == nil {
		t.Fatal("Client should not be nil")
	}
}

func TestNewClientWithCustomAuth(t *testing.T) {
	customAuth := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("X-Custom-Token", "test-token")
		return nil
	}
	
	client, err := NewClientWithCustomAuth("https://api.example.com", customAuth)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	
	if client == nil {
		t.Fatal("Client should not be nil")
	}
}

func TestClient_V1(t *testing.T) {
	client, err := NewClient("https://api.example.com", "test-api-key")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	
	v1 := client.V1()
	if v1 == nil {
		t.Fatal("V1API should not be nil")
	}
}

func TestNewMeteringRecord(t *testing.T) {
	startTime := time.Now()
	record := NewMeteringRecord("product-123", "customer-456", "api-requests", 100, startTime)
	
	if record.ProductID != "product-123" {
		t.Errorf("expected ProductID 'product-123', got '%s'", record.ProductID)
	}
	
	if record.CustomerID != "customer-456" {
		t.Errorf("expected CustomerID 'customer-456', got '%s'", record.CustomerID)
	}
	
	if record.Dimension.Name != "api-requests" {
		t.Errorf("expected dimension name 'api-requests', got '%s'", record.Dimension.Name)
	}
	
	if record.Dimension.Quantity != 100 {
		t.Errorf("expected quantity 100, got %f", record.Dimension.Quantity)
	}
	
	if record.StartTime != startTime.Format(time.RFC3339) {
		t.Errorf("expected start time '%s', got '%s'", startTime.Format(time.RFC3339), record.StartTime)
	}
}

func TestNewMeteringRecordWithCloud(t *testing.T) {
	startTime := time.Now()
	record := NewMeteringRecordWithCloud("aws", "product-123", "customer-456", "compute-hours", 50, startTime)
	
	if record.Cloud != "aws" {
		t.Errorf("expected cloud 'aws', got '%s'", record.Cloud)
	}
	
	if record.ProductID != "product-123" {
		t.Errorf("expected ProductID 'product-123', got '%s'", record.ProductID)
	}
	
	if record.Dimension.Name != "compute-hours" {
		t.Errorf("expected dimension name 'compute-hours', got '%s'", record.Dimension.Name)
	}
	
	if record.Dimension.Quantity != 50 {
		t.Errorf("expected quantity 50, got %f", record.Dimension.Quantity)
	}
}

func TestMeteringRecord_WithCloud(t *testing.T) {
	record := NewMeteringRecord("product-123", "customer-456", "api-requests", 100, time.Now())
	recordWithCloud := record.WithCloud("azure")
	
	if recordWithCloud.Cloud != "azure" {
		t.Errorf("expected cloud 'azure', got '%s'", recordWithCloud.Cloud)
	}
	
	// Original record should be unchanged
	if record.Cloud != "" {
		t.Errorf("original record cloud should be empty, got '%s'", record.Cloud)
	}
}

func TestMeteringRecord_WithScheduledAt(t *testing.T) {
	scheduledTime := time.Now().Add(1 * time.Hour)
	record := NewMeteringRecord("product-123", "customer-456", "api-requests", 100, time.Now())
	recordWithSchedule := record.WithScheduledAt(scheduledTime)
	
	if recordWithSchedule.ScheduledAt != scheduledTime.Format(time.RFC3339) {
		t.Errorf("expected scheduled time '%s', got '%s'", scheduledTime.Format(time.RFC3339), recordWithSchedule.ScheduledAt)
	}
	
	// Original record should be unchanged
	if record.ScheduledAt != "" {
		t.Errorf("original record scheduled time should be empty, got '%s'", record.ScheduledAt)
	}
}

func TestMeteringRecordBuilder(t *testing.T) {
	startTime := time.Now()
	scheduledTime := startTime.Add(1 * time.Hour)
	
	record := NewMeteringRecordBuilder().
		ProductID("product-123").
		CustomerID("customer-456").
		Dimension("api-requests", 100).
		StartTime(startTime).
		Cloud("aws").
		ScheduledAt(scheduledTime).
		Build()
	
	if record.ProductID != "product-123" {
		t.Errorf("expected ProductID 'product-123', got '%s'", record.ProductID)
	}
	
	if record.CustomerID != "customer-456" {
		t.Errorf("expected CustomerID 'customer-456', got '%s'", record.CustomerID)
	}
	
	if record.Dimension.Name != "api-requests" {
		t.Errorf("expected dimension name 'api-requests', got '%s'", record.Dimension.Name)
	}
	
	if record.Dimension.Quantity != 100 {
		t.Errorf("expected quantity 100, got %f", record.Dimension.Quantity)
	}
	
	if record.StartTime != startTime.Format(time.RFC3339) {
		t.Errorf("expected start time '%s', got '%s'", startTime.Format(time.RFC3339), record.StartTime)
	}
	
	if record.Cloud != "aws" {
		t.Errorf("expected cloud 'aws', got '%s'", record.Cloud)
	}
	
	if record.ScheduledAt != scheduledTime.Format(time.RFC3339) {
		t.Errorf("expected scheduled time '%s', got '%s'", scheduledTime.Format(time.RFC3339), record.ScheduledAt)
	}
}

func TestNewUsageRecord(t *testing.T) {
	startTime := time.Now()
	record := NewUsageRecord("product-123", "customer-456", "api-requests", 100, startTime)
	
	if record.ProductId != "product-123" {
		t.Errorf("expected ProductId 'product-123', got '%s'", record.ProductId)
	}
	
	if record.CustomerIdentifier != "customer-456" {
		t.Errorf("expected CustomerIdentifier 'customer-456', got '%s'", record.CustomerIdentifier)
	}
	
	if record.Dimension.Name != "api-requests" {
		t.Errorf("expected dimension name 'api-requests', got '%s'", record.Dimension.Name)
	}
	
	if record.Dimension.Quantity != 100 {
		t.Errorf("expected quantity 100, got %f", record.Dimension.Quantity)
	}
	
	if !record.StartTime.Equal(startTime) {
		t.Errorf("expected start time %v, got %v", startTime, record.StartTime)
	}
}

func TestNewUsageAllocation(t *testing.T) {
	tags := []usageapi.Tag{
		NewTag("environment", "production"),
		NewTag("team", "backend"),
	}
	
	allocation := NewUsageAllocation(50, tags)
	
	if allocation.AllocatedUsageQuantity == nil {
		t.Fatal("AllocatedUsageQuantity should not be nil")
	}
	
	if *allocation.AllocatedUsageQuantity != 50 {
		t.Errorf("expected quantity 50, got %f", *allocation.AllocatedUsageQuantity)
	}
	
	if allocation.Tags == nil {
		t.Fatal("Tags should not be nil")
	}
	
	if len(*allocation.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(*allocation.Tags))
	}
}

func TestNewTag(t *testing.T) {
	tag := NewTag("environment", "production")
	
	if tag.Key != "environment" {
		t.Errorf("expected key 'environment', got '%s'", tag.Key)
	}
	
	if tag.Value != "production" {
		t.Errorf("expected value 'production', got '%s'", tag.Value)
	}
}