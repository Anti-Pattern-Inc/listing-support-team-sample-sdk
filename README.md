# AWS Marketplace SaaS Usage Records SDK for Go

AWS Marketplace SaaSサービス向けの使用量レコードAPI用Go SDKです。OpenAPIから自動生成されたクライアントとシンプルなラッパーを提供します。

## 特徴

- 🚀 **型安全**: OpenAPIから自動生成されたGoクライアント
- 🔐 **認証対応**: 環境変数から自動認証
- 🛠️ **シンプルなAPI**: SaaSus SDKスタイルの設計
- 📦 **エラーハンドリング**: 詳細なエラーレスポンス対応
- 🎯 **使いやすい**: 認証ボイラープレートを排除した設計

## インストール

```bash
go get github.com/Anti-Pattern-Inc/listing-support-team-sample-sdk
```

## クイックスタート

### 環境変数設定

```bash
export MARKETPLACE_CLIENT_ID="your_client_id"
export MARKETPLACE_API_KEY="your_api_key"
```

### 基本的な使用方法

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

    // 1. 環境変数から自動認証クライアント作成
    client, err := marketplace.ClientWithResponse(ctx)
    if err != nil {
        log.Fatal("Failed to create client:", err)
    }

    // 2. 使用量レコードの作成
    records := []clientapi.UsageRecord{
        {
            ProductId:          "my-product-id",
            CustomerIdentifier: "customer-123",
            Dimension: clientapi.Dimension{
                Name:     "api-calls",
                Quantity: 100,
            },
            StartTime: time.Now(),
        },
    }

    // 3. APIコール
    response, err := client.CreateUsageRecordsWithResponse(ctx, records)
    if err != nil {
        log.Fatal("API call failed:", err)
    }

    // 4. レスポンス処理
    switch response.StatusCode() {
    case 200:
        fmt.Println("✅ Success")
        if response.JSON200 != nil && response.JSON200.TotalRecords != nil {
            fmt.Printf("Total records: %.0f\n", *response.JSON200.TotalRecords)
        }
    case 207:
        fmt.Println("⚠️ Partial success")
    case 401:
        fmt.Printf("❌ Authentication error: %s\n", response.JSON401.Error)
    case 400:
        fmt.Printf("❌ Bad request: %s\n", response.JSON400.Error)
    case 500:
        fmt.Printf("❌ Server error: %s\n", response.JSON500.Error)
    default:
        fmt.Printf("❌ Unexpected status: %d\n", response.StatusCode())
    }
}
```

## API使用方法

### クライアント作成

```go
// 環境変数 MARKETPLACE_CLIENT_ID と MARKETPLACE_API_KEY から自動認証
client, err := marketplace.ClientWithResponse(ctx)
```

### 利用可能なAPI

SaaSus SDKスタイルでクライアント取得後、直接APIを呼び出せます：

- `client.CreateUsageRecordsWithResponse(ctx, records)`: 使用量記録作成
- `client.UpdateUsageRecordsWithResponse(ctx, records)`: 使用量記録更新
- `client.GetUsageRecordsWithResponse(ctx, params)`: 使用量記録取得
- `client.OptionsUsageRecordsWithResponse(ctx)`: CORS対応
- `client.IssueAuthTokenWithResponse(ctx, request)`: 認証トークン取得

### 使用量レコード取得

```go
func getUsageRecords() error {
    ctx := context.Background()
    
    // 環境変数から自動認証
    client, err := marketplace.ClientWithResponse(ctx)
    if err != nil {
        return err
    }

    // 期間指定で取得（過去24時間）
    now := time.Now()
    yesterday := now.AddDate(0, 0, -1)
    
    allStatus := clientapi.All
    params := &clientapi.GetUsageRecordsParams{
        ProductId:      stringPtr("my-product-id"),
        StartTimestamp: stringPtr(fmt.Sprintf("%d", yesterday.Unix())),
        EndTimestamp:   stringPtr(fmt.Sprintf("%d", now.Unix())),
        Limit:          stringPtr("10"),
        Status:         &allStatus,
    }

    response, err := client.GetUsageRecordsWithResponse(ctx, params)
    if err != nil {
        return err
    }
}

func stringPtr(s string) *string { return &s }
```

## 使用量レコードの構造

### 基本的なレコード

```go
record := clientapi.UsageRecord{
    ProductId:          "product-123",           // 製品ID
    CustomerIdentifier: "customer-456",         // 顧客識別子
    Dimension: clientapi.Dimension{
        Name:     "api-requests",               // メトリクス名
        Quantity: 100,                         // 使用量
    },
    StartTime: time.Now(),                     // 開始時刻
}
```

### 使用量アロケーション付き

```go
record := clientapi.UsageRecord{
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
        },
    },
    StartTime: time.Now(),
}

func floatPtr(f float64) *float64 { return &f }
```

## エラーハンドリング

```go
response, err := client.CreateUsageRecordsWithResponse(ctx, records)
if err != nil {
    // ネットワークエラーやJSON解析エラー
    log.Printf("Request error: %v", err)
    return
}

// HTTPステータスコードによる処理
switch response.StatusCode() {
case 200:
    // 全レコード成功
    fmt.Println("All records processed successfully")
case 207:
    // 部分的成功
    if response.JSON207 != nil {
        fmt.Printf("Success: %.0f, Failed: %.0f\n", 
            *response.JSON207.SuccessfulRecords, 
            *response.JSON207.FailedRecords)
    }
case 400:
    // リクエストエラー
    fmt.Printf("Bad request: %s\n", response.JSON400.Error)
case 401:
    // 認証エラー
    fmt.Printf("Authentication failed: %s\n", response.JSON401.Error)
case 409:
    // 更新時の競合エラー
    fmt.Printf("Conflict: %s\n", response.JSON409.Error)
case 500:
    // サーバーエラー
    fmt.Printf("Server error: %s\n", response.JSON500.Error)
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
├── middleware/
│   └── authenticate.go         # 認証ミドルウェア
└── modules/
    └── clientapi/
        └── client.go           # SaaSusスタイルクライアント
```

## SDK構造

### modules/clientapi（推奨）

SaaSus SDKスタイルの使いやすいクライアント：

```go
import marketplace "github.com/Anti-Pattern-Inc/listing-support-team-sample-sdk/modules/clientapi"

// 唯一のAPIエントリーポイント
func ClientWithResponse(ctx context.Context) (*clientapi.ClientWithResponses, error)
```

### generated/clientapi（低レベルAPI）

OpenAPIから自動生成されたクライアント：

```go
import "github.com/Anti-Pattern-Inc/listing-support-team-sample-sdk/generated/clientapi"

// 直接使用する場合
client, err := clientapi.NewClientWithResponses(baseURL)
```

### middleware

認証処理：

```go
import "github.com/Anti-Pattern-Inc/listing-support-team-sample-sdk/middleware"

// 認証トークン取得（内部で使用）
token, err := middleware.Authenticate(ctx, clientID, apiKey)
```

## 実用的な使用例

### バッチ処理での使用量記録

```go
func recordUsageBatch(usageData []UsageData) error {
    ctx := context.Background()
    
    // 環境変数から認証
    client, err := marketplace.ClientWithResponse(ctx)
    if err != nil {
        return fmt.Errorf("failed to create client: %w", err)
    }

    var records []clientapi.UsageRecord
    for _, usage := range usageData {
        records = append(records, clientapi.UsageRecord{
            ProductId:          usage.ProductID,
            CustomerIdentifier: usage.CustomerID,
            Dimension: clientapi.Dimension{
                Name:     usage.MetricName,
                Quantity: usage.Quantity,
            },
            StartTime: usage.Timestamp,
        })
    }

    response, err := client.CreateUsageRecordsWithResponse(ctx, records)
    if err != nil {
        return fmt.Errorf("API call failed: %w", err)
    }

    if response.StatusCode() != 200 && response.StatusCode() != 207 {
        return fmt.Errorf("API returned error status: %d", response.StatusCode())
    }

    return nil
}
```

### Webアプリケーションでの使用

```go
func usageHandler(w http.ResponseWriter, r *http.Request) {
    // リクエストから使用量データを取得
    var requestData struct {
        ProductID  string  `json:"product_id"`
        CustomerID string  `json:"customer_id"`
        Metric     string  `json:"metric"`
        Quantity   float64 `json:"quantity"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    // クライアント作成（環境変数から認証）
    client, err := marketplace.ClientWithResponse(r.Context())
    if err != nil {
        http.Error(w, "Authentication failed", http.StatusUnauthorized)
        return
    }

    // 使用量記録作成
    records := []clientapi.UsageRecord{
        {
            ProductId:          requestData.ProductID,
            CustomerIdentifier: requestData.CustomerID,
            Dimension: clientapi.Dimension{
                Name:     requestData.Metric,
                Quantity: requestData.Quantity,
            },
            StartTime: time.Now(),
        },
    }

    response, err := client.CreateUsageRecordsWithResponse(r.Context(), records)
    if err != nil {
        http.Error(w, "API call failed", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": response.StatusCode() == 200,
        "status":  response.StatusCode(),
    })
}
```

### シンプルな使用例

```go
func createUsageRecord() error {
    ctx := context.Background()
    
    // 1行でクライアント取得
    client, err := marketplace.ClientWithResponse(ctx)
    if err != nil {
        return err
    }

    // レコード作成とAPI呼び出し
    records := []clientapi.UsageRecord{
        {
            ProductId:          "my-product",
            CustomerIdentifier: "customer-123",
            Dimension: clientapi.Dimension{
                Name:     "api-calls",
                Quantity: 1,
            },
            StartTime: time.Now(),
        },
    }

    response, err := client.CreateUsageRecordsWithResponse(ctx, records)
    if err != nil {
        return err
    }

    fmt.Printf("Status: %d\n", response.StatusCode())
    return nil
}
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

## トラブルシューティング

### 認証エラー

```
❌ Authentication error: invalid credentials
```

- `MARKETPLACE_CLIENT_ID`と`MARKETPLACE_API_KEY`の値を確認
- 環境変数が正しく設定されているか確認

### 400 Bad Request

```
❌ Bad request: invalid usage record format
```

- `ProductId`、`CustomerIdentifier`、`Dimension`が正しく設定されているか確認
- `StartTime`が適切な形式か確認

### 409 Conflict (Update時)

```
❌ Conflict: records not updatable
```

- レコードが更新可能な状態か確認
- 同じレコードを重複して更新していないか確認

## ライセンス

[ライセンス情報をここに記載]