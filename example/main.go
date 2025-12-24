package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Anti-Pattern-Inc/listing-support-team-sample-sdk/generated/clientapi"
	marketplace "github.com/Anti-Pattern-Inc/listing-support-team-sample-sdk/modules/clientapi"
)

func main() {
	ctx := context.Background()

	// Create client from environment variables (CLIENT_ID and API_KEY)
	client, err := marketplace.ClientWithResponse(ctx)
	if err != nil {
		log.Fatal("Failed to create client:", err)
	}

	// Create sample usage records
	records := []clientapi.UsageRecord{
		{
			ProductId:          "sample-product-id",
			CustomerIdentifier: "customer-123",
			Dimension: clientapi.Dimension{
				Name:     "api-calls",
				Quantity: 100,
			},
			StartTime: time.Now(),
		},
		{
			ProductId:          "sample-product-id",
			CustomerIdentifier: "customer-456",
			Dimension: clientapi.Dimension{
				Name:     "storage-gb",
				Quantity: 50,
			},
			StartTime: time.Now(),
		},
	}

	fmt.Println("Sending usage records to AWS Marketplace...")

	// Call API
	response, err := client.CreateUsageRecordsWithResponse(ctx, records)
	if err != nil {
		log.Fatal("API call failed:", err)
	}

	// Handle response based on status code
	switch response.StatusCode() {
	case 200:
		fmt.Println("✅ Success: All records processed successfully")
		if response.JSON200 != nil && response.JSON200.TotalRecords != nil {
			fmt.Printf("Total records: %.0f\n", *response.JSON200.TotalRecords)
		}
	case 207:
		fmt.Println("⚠️ Partial success: Some records failed")
		if response.JSON207 != nil {
			if response.JSON207.SuccessfulRecords != nil {
				fmt.Printf("Successful: %.0f\n", *response.JSON207.SuccessfulRecords)
			}
			if response.JSON207.FailedRecords != nil {
				fmt.Printf("Failed: %.0f\n", *response.JSON207.FailedRecords)
			}
		}
	case 400:
		fmt.Printf("❌ Bad request: %s\n", response.JSON400.Error)
	case 401:
		fmt.Printf("❌ Authentication error: %s\n", response.JSON401.Error)
	case 500:
		fmt.Printf("❌ Server error: %s\n", response.JSON500.Error)
	default:
		fmt.Printf("❌ Unexpected status: %d\n", response.StatusCode())
	}
}
