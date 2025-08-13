package usage

import (
	"time"

	"github.com/nishimurashinya/openapi-go-sdk/generated/usageapi"
)

// MeteringRecord represents a simplified usage record for metering (Node.js-style)
type MeteringRecord struct {
	// Cloud platform (e.g., "aws", "azure", "gcp")
	Cloud string `json:"cloud,omitempty"`
	// ProductID is the marketplace product identifier
	ProductID string `json:"product_id"`
	// CustomerID is the customer identifier
	CustomerID string `json:"customer_id"`
	// Dimension contains the usage dimension information
	Dimension MeteringDimension `json:"dimension"`
	// StartTime is when the usage started (ISO 8601 format)
	StartTime string `json:"start_time"`
	// ScheduledAt is when the record should be processed (optional)
	ScheduledAt string `json:"scheduled_at,omitempty"`
}

// MeteringDimension represents usage dimension in a simplified format
type MeteringDimension struct {
	// Name of the dimension (e.g., "api-requests", "compute-hours")
	Name string `json:"name"`
	// Quantity of usage
	Quantity float32 `json:"quantity"`
}

// NewMeteringRecord creates a new MeteringRecord with required fields
func NewMeteringRecord(productID, customerID, dimensionName string, quantity float32, startTime time.Time) MeteringRecord {
	return MeteringRecord{
		ProductID:  productID,
		CustomerID: customerID,
		Dimension: MeteringDimension{
			Name:     dimensionName,
			Quantity: quantity,
		},
		StartTime: startTime.Format(time.RFC3339),
	}
}

// NewMeteringRecordWithCloud creates a new MeteringRecord with cloud platform specified
func NewMeteringRecordWithCloud(cloud, productID, customerID, dimensionName string, quantity float32, startTime time.Time) MeteringRecord {
	record := NewMeteringRecord(productID, customerID, dimensionName, quantity, startTime)
	record.Cloud = cloud
	return record
}

// NewMeteringRecordWithSchedule creates a new MeteringRecord with scheduled processing time
func NewMeteringRecordWithSchedule(productID, customerID, dimensionName string, quantity float32, startTime, scheduledAt time.Time) MeteringRecord {
	record := NewMeteringRecord(productID, customerID, dimensionName, quantity, startTime)
	record.ScheduledAt = scheduledAt.Format(time.RFC3339)
	return record
}

// WithCloud sets the cloud platform for the record (method chaining)
func (r MeteringRecord) WithCloud(cloud string) MeteringRecord {
	r.Cloud = cloud
	return r
}

// WithScheduledAt sets the scheduled processing time (method chaining)
func (r MeteringRecord) WithScheduledAt(scheduledAt time.Time) MeteringRecord {
	r.ScheduledAt = scheduledAt.Format(time.RFC3339)
	return r
}

// MeteringRecordBuilder provides a builder pattern for creating MeteringRecord
type MeteringRecordBuilder struct {
	record MeteringRecord
}

// NewMeteringRecordBuilder creates a new builder for MeteringRecord
func NewMeteringRecordBuilder() *MeteringRecordBuilder {
	return &MeteringRecordBuilder{}
}

// ProductID sets the product ID
func (b *MeteringRecordBuilder) ProductID(productID string) *MeteringRecordBuilder {
	b.record.ProductID = productID
	return b
}

// CustomerID sets the customer ID
func (b *MeteringRecordBuilder) CustomerID(customerID string) *MeteringRecordBuilder {
	b.record.CustomerID = customerID
	return b
}

// Dimension sets the usage dimension
func (b *MeteringRecordBuilder) Dimension(name string, quantity float32) *MeteringRecordBuilder {
	b.record.Dimension = MeteringDimension{
		Name:     name,
		Quantity: quantity,
	}
	return b
}

// StartTime sets the start time
func (b *MeteringRecordBuilder) StartTime(startTime time.Time) *MeteringRecordBuilder {
	b.record.StartTime = startTime.Format(time.RFC3339)
	return b
}

// Cloud sets the cloud platform
func (b *MeteringRecordBuilder) Cloud(cloud string) *MeteringRecordBuilder {
	b.record.Cloud = cloud
	return b
}

// ScheduledAt sets the scheduled processing time
func (b *MeteringRecordBuilder) ScheduledAt(scheduledAt time.Time) *MeteringRecordBuilder {
	b.record.ScheduledAt = scheduledAt.Format(time.RFC3339)
	return b
}

// Build returns the constructed MeteringRecord
func (b *MeteringRecordBuilder) Build() MeteringRecord {
	return b.record
}

// Legacy helper functions for backward compatibility with low-level API

// NewUsageRecord creates a new UsageRecord with required fields (low-level API)
func NewUsageRecord(productID, customerID, dimensionName string, quantity float32, startTime time.Time) usageapi.UsageRecord {
	return usageapi.UsageRecord{
		ProductId:          productID,
		CustomerIdentifier: customerID,
		Dimension: usageapi.Dimension{
			Name:     dimensionName,
			Quantity: quantity,
		},
		StartTime: startTime,
	}
}

// NewUsageRecordWithAllocations creates a new UsageRecord with usage allocations (low-level API)
func NewUsageRecordWithAllocations(productID, customerID, dimensionName string, quantity float32, startTime time.Time, allocations []usageapi.UsageAllocation) usageapi.UsageRecord {
	record := NewUsageRecord(productID, customerID, dimensionName, quantity, startTime)
	record.Dimension.UsageAllocations = &allocations
	return record
}

// NewUsageAllocation creates a new UsageAllocation (low-level API)
func NewUsageAllocation(quantity float32, tags []usageapi.Tag) usageapi.UsageAllocation {
	return usageapi.UsageAllocation{
		AllocatedUsageQuantity: &quantity,
		Tags:                   &tags,
	}
}

// NewTag creates a new Tag (low-level API)
func NewTag(key, value string) usageapi.Tag {
	return usageapi.Tag{
		Key:   key,
		Value: value,
	}
}