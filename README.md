# AWS Marketplace SaaS Usage Records SDK for Go

AWS Marketplace SaaSサービス向けの使用量レコードAPI用Go SDKです。OpenAPIから自動生成されたクライアントとユーザーフレンドリーなラッパーを提供します。

## 特徴

- 🚀 **型安全**: OpenAPIから自動生成されたGoクライアント
- 🔐 **認証対応**: ClientID + APIKey による認証
- 🛠️ **シンプルなAPI**: SaaSusスタイルのクライアントインターフェース
- 📦 **エラーハンドリング**: 詳細なエラーレスポンス対応
- 🎯 **使いやすい**: 認証ボイラープレートを排除した設計

## インストール

```bash
go get github.com/Anti-Pattern-Inc/listing-support-team-sample-sdk
```

## クイックスタート

### SaaSusスタイルAPI（推奨）

```go
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

    fmt.Println("🚀 SaaSus-style API client example")

    // Example 1: Simple authentication + marketplace operations (like SaaSus SDK)
    fmt.Println("\n📝 Example 1: SaaSus-style usage")
    
    clientID := "your_client_id"
    apiKey := "your_api_key"

    // Step 1: Create marketplace client (handles auth internally)
    marketplaceClient, err := marketplace.NewMarketplaceClientWithAuth(ctx, clientID, apiKey)
    if err != nil {
        log.Fatal("❌ Failed to create marketplace client:", err)
    }
    fmt.Println("✅ Marketplace client created")

    // Step 2: Create usage records (simple method call)
    records := []clientapi.UsageRecord{
        {
            ProductId:          "my-product",
            CustomerIdentifier: "customer-123",
            Dimension: clientapi.Dimension{
                Name:     "api-calls",
                Quantity: 100,
            },
            StartTime: time.Now(),
        },
    }

    response, err := marketplaceClient.Client.CreateUsageRecordsWithResponse(ctx, records)
    if err != nil {
        log.Printf("❌ Failed to create records: %v", err)
    } else {
        handleResponse("Create", response)
    }

    // Example 2: Separate auth and marketplace clients (more control)
    fmt.Println("\n📝 Example 2: Separate auth client")
    
    // Create auth client
    authClient, err := marketplace.NewAuthClient()
    if err != nil {
        log.Fatal("❌ Failed to create auth client:", err)
    }

    // Get token
    token, err := authClient.GetAuthToken(ctx, clientID, apiKey)
    if err != nil {
        log.Printf("❌ Authentication failed: %v", err)
        return
    }
    fmt.Println("✅ Token obtained")

    // Create marketplace client with token
    marketplaceClient2, err := marketplace.NewMarketplaceClient(token)
    if err != nil {
        log.Fatal("❌ Failed to create marketplace client:", err)
    }

    // Update usage records
    updateRecords := []clientapi.UsageRecord{
        {
            ProductId:          "my-product",
            CustomerIdentifier: "customer-456",
            Dimension: clientapi.Dimension{
                Name:     "storage-gb",
                Quantity: 50.5,
            },
            StartTime: time.Now(),
        },
    }

    updateResponse, err := marketplaceClient2.Client.UpdateUsageRecordsWithResponse(ctx, updateRecords)
    if err != nil {
        log.Printf("❌ Failed to update records: %v", err)
    } else {
        handleUpdateResponse("Update", updateResponse)
    }

    fmt.Println("\n🎉 SaaSus-style examples completed!")
}

func handleResponse(operation string, response *clientapi.CreateUsageRecordsResponse) {
    switch response.StatusCode() {
    case 200:
        fmt.Printf("✅ %s: Success\n", operation)
        if response.JSON200 != nil && response.JSON200.TotalRecords != nil {
            fmt.Printf("   Total records: %.0f\n", *response.JSON200.TotalRecords)
        }
    case 207:
        fmt.Printf("⚠️ %s: Partial success\n", operation)
        if response.JSON207 != nil {
            if response.JSON207.SuccessfulRecords != nil && response.JSON207.FailedRecords != nil {
                fmt.Printf("   Success: %.0f, Failed: %.0f\n", 
                    *response.JSON207.SuccessfulRecords, *response.JSON207.FailedRecords)
            }
        }
    case 401:
        fmt.Printf("❌ %s: Authentication error: %s\n", operation, response.JSON401.Error)
    case 400:
        fmt.Printf("❌ %s: Bad request: %s\n", operation, response.JSON400.Error)
    case 500:
        fmt.Printf("❌ %s: Server error: %s\n", operation, response.JSON500.Error)
    default:
        fmt.Printf("❌ %s: Unexpected status: %d\n", operation, response.StatusCode())
    }
}

func handleUpdateResponse(operation string, response *clientapi.UpdateUsageRecordsResponse) {
    switch response.StatusCode() {
    case 200:
        fmt.Printf("✅ %s: Success\n", operation)
        if response.JSON200 != nil && response.JSON200.TotalRecords != nil {
            fmt.Printf("   Total records: %.0f\n", *response.JSON200.TotalRecords)
        }
    case 207:
        fmt.Printf("⚠️ %s: Partial success\n", operation)
        if response.JSON207 != nil {
            if response.JSON207.SuccessfulRecords != nil && response.JSON207.FailedRecords != nil {
                fmt.Printf("   Success: %.0f, Failed: %.0f\n", 
                    *response.JSON207.SuccessfulRecords, *response.JSON207.FailedRecords)
            }
        }
    case 401:
        fmt.Printf("❌ %s: Authentication error: %s\n", operation, response.JSON401.Error)
    case 400:
        fmt.Printf("❌ %s: Bad request: %s\n", operation, response.JSON400.Error)
    case 409:
        fmt.Printf("⚠️ %s: Conflict - records not updatable: %s\n", operation, response.JSON409.Error)
    case 500:
        fmt.Printf("❌ %s: Server error: %s\n", operation, response.JSON500.Error)
    default:
        fmt.Printf("❌ %s: Unexpected status: %d\n", operation, response.StatusCode())
    }
}
```

### シンプルなワンライナー

```go
// 最もシンプルな使用方法
client, err := marketplace.NewMarketplaceClientWithAuth(ctx, clientID, apiKey)
response, err := client.Client.CreateUsageRecordsWithResponse(ctx, records)
```

## API構造

### modules/clientapi（高レベルAPI）- 推奨

SaaSusスタイルのシンプルなクライアントインターフェース：

```go
import marketplace "github.com/Anti-Pattern-Inc/listing-support-team-sample-sdk/modules/clientapi"

// 認証込みワンライナー
client, err := marketplace.NewMarketplaceClientWithAuth(ctx, clientID, apiKey)

// または段階的
authClient, err := marketplace.NewAuthClient()
token, err := authClient.GetAuthToken(ctx, clientID, apiKey)
marketplaceClient, err := marketplace.NewMarketplaceClient(token)
```

### generated/clientapi（低レベルAPI）

OpenAPIから自動生成されたクライアントを直接使用：

```go
import "github.com/Anti-Pattern-Inc/listing-support-team-sample-sdk/generated/clientapi"

// 基本クライアント
client, err := clientapi.NewClientWithResponses(baseURL)

// 認証付きクライアント
client, err := clientapi.NewClientWithResponses(
    baseURL,
    clientapi.WithRequestEditorFn(authFunc),
)
```

## 利用可能なAPI

### 認証
- `GetAuthToken(ctx, clientID, apiKey)`: 認証トークン取得

### 使用量記録
- `Client.CreateUsageRecordsWithResponse(ctx, records)`: 使用量記録作成
- `Client.UpdateUsageRecordsWithResponse(ctx, records)`: 使用量記録更新
- その他すべてのAPIメソッドが`Client`経由で利用可能

### 低レベルAPI
- `IssueAuthTokenWithResponse()`: 認証トークン取得
- `CreateUsageRecordsWithResponse()`: 使用量記録作成
- `UpdateUsageRecordsWithResponse()`: 使用量記録更新
- `OptionsUsageRecordsWithResponse()`: CORS対応

## 使用量レコードの作成

### 基本的な使用量レコード

```go
records := []clientapi.UsageRecord{
    {
        ProductId:          "product-123",
        CustomerIdentifier: "customer-456",
        Dimension: clientapi.Dimension{
            Name:     "api-requests",
            Quantity: 100,
        },
        StartTime: time.Now(),
    },
}

response, err := client.Client.CreateUsageRecordsWithResponse(ctx, records)
```

### 使用量アロケーション付き

```go
records := []clientapi.UsageRecord{
    {
        ProductId:          "product-storage",
        CustomerIdentifier: "customer-789",
        Dimension: clientapi.Dimension{
            Name:     "storage-gb",
            Quantity: 50,
            UsageAllocations: &[]clientapi.UsageAllocation{
                {
                    AllocatedUsageQuantity: floatPtr(25),
                    Tags: &[]clientapi.Tag{
                        {Key: "region", Value: "us-east-1"},
                        {Key: "tier", Value: "premium"},
                    },
                },
                {
                    AllocatedUsageQuantity: floatPtr(25),
                    Tags: &[]clientapi.Tag{
                        {Key: "region", Value: "us-west-2"},
                        {Key: "tier", Value: "standard"},
                    },
                },
            },
        },
        StartTime: time.Now(),
    },
}
```

## エラーハンドリング

```go
response, err := client.Client.CreateUsageRecordsWithResponse(ctx, records)
if err != nil {
    log.Printf("API Error: %v", err)
    return
}

switch response.StatusCode() {
case 200:
    fmt.Printf("✅ Success\n")
    if response.JSON200 != nil && response.JSON200.TotalRecords != nil {
        fmt.Printf("   Total records: %.0f\n", *response.JSON200.TotalRecords)
    }
case 207:
    fmt.Printf("⚠️ Partial Success\n")
    if response.JSON207 != nil {
        if response.JSON207.SuccessfulRecords != nil && response.JSON207.FailedRecords != nil {
            fmt.Printf("   Success: %.0f, Failed: %.0f\n", 
                *response.JSON207.SuccessfulRecords, *response.JSON207.FailedRecords)
        }
    }
case 400:
    fmt.Printf("❌ Bad Request: %s\n", response.JSON400.Error)
case 401:
    fmt.Printf("❌ Unauthorized: %s\n", response.JSON401.Error)
case 409:
    fmt.Printf("❌ Conflict: %s\n", response.JSON409.Error)
case 500:
    fmt.Printf("❌ Server Error: %s\n", response.JSON500.Error)
default:
    fmt.Printf("❌ Unexpected status: %d\n", response.StatusCode())
}
```

## ファイル構造

```
├── README.md                    # このファイル
├── CLAUDE.md                    # 開発者向けガイド
├── go.mod                       # Goモジュール定義
├── openapi.yaml                 # OpenAPI仕様書
├── generate.sh                  # クライアント生成スクリプト
├── generated/
│   └── clientapi/
│       ├── client.gen.go       # 生成されたクライアント
│       └── types.gen.go        # 生成された型定義
└── modules/
    └── clientapi/
        ├── auth.go             # 認証クライアント
        └── client.go           # マーケットプレイスクライアント
```

## 開発

### クライアント再生成

```bash
./generate.sh
```

### テスト実行

```bash
go test -v ./...
```

### ビルド

```bash
go build ./...
```

## サンプルコード

完全なサンプルコードはREADME冒頭の「クイックスタート」セクションを参照してください。

## ライセンス

[ライセンス情報をここに記載]

## 貢献

プルリクエストやイシューの報告を歓迎します。

## サポート

- [OpenAPI仕様](./openapi.yaml)
- [開発者向けドキュメント](./CLAUDE.md)