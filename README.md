# AWS Marketplace SaaS Usage Records SDK for Go

AWS Marketplace SaaSサービス向けの使用量レコードAPI用Go SDKです。Node.js SDKと同様の使いやすいインターフェースを提供します。

## 特徴

- 🚀 **Node.js風のシンプルAPI**: `client.V1().CreateMeteringUsageRecords()`
- 🔐 **複数の認証方式**: Bearer、HMAC-SHA256署名、カスタム認証
- 🛠️ **便利なヘルパー関数**: メソッドチェーン、ビルダーパターン対応
- 📦 **型安全**: GoのStructによる型安全性
- 🔧 **環境変数サポート**: `AWS_MARKETPLACE_API_KEY`, `AWS_MARKETPLACE_SECRET_KEY`

## インストール

```bash
go get github.com/Anti-Pattern-Inc/listing-support-team-sample-sdk
```

## クイックスタート

### 基本的な使用方法

```go
package main

import (
    "context"
    "log"
    "time"
    
    "github.com/Anti-Pattern-Inc/listing-support-team-sample-sdk/modules/usage"
)

func main() {
    // Usage APIクライアント作成（Base URLを直接指定）
    client, err := usage.NewClient(
        "https://api-id.execute-api.ap-northeast-1.amazonaws.com/Prod",
        "your-api-key",
    )
    if err != nil {
        log.Fatal(err)
    }

    // 使用量レコード作成（Node.js風）
    records := []usage.MeteringRecord{
        {
            Cloud:      "aws",
            ProductID:  "1test1a2b3c4d5e6f7g8h9ijk",
            CustomerID: "ij3sXMkN3or",
            Dimension:  usage.MeteringDimension{Name: "api-requests", Quantity: 100},
            StartTime:  time.Now().Format(time.RFC3339),
        },
    }

    ctx := context.Background()
    response, err := client.V1().CreateMeteringUsageRecords(ctx, records)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Status: %d", response.StatusCode())
}
```

## 認証方式

### 1. Bearer認証（推奨）

```go
client, err := usage.NewClient(
    "https://api-id.execute-api.ap-northeast-1.amazonaws.com/Prod", 
    "your-api-key",
)
```

### 2. HMAC-SHA256署名認証（高セキュリティ）

```go
client, err := usage.NewClientWithSignature(
    "https://api-id.execute-api.ap-northeast-1.amazonaws.com/Prod",
    "api-key", 
    "secret-key",
)
```

### 3. カスタム認証

```go
customAuth := func(ctx context.Context, req *http.Request) error {
    req.Header.Set("X-Custom-Token", "your-token")
    return nil
}

client, err := usage.NewClientWithCustomAuth(
    "https://api-id.execute-api.ap-northeast-1.amazonaws.com/Prod",
    customAuth,
)
```

## 使用量レコードの作成方法

### 基本的な作成

```go
record := usage.NewMeteringRecord(
    "product-id", 
    "customer-id", 
    "api-requests", 
    100, 
    time.Now(),
)
```

### クラウドプラットフォーム指定

```go
record := usage.NewMeteringRecordWithCloud(
    "aws", 
    "product-id", 
    "customer-id", 
    "compute-hours", 
    50, 
    time.Now(),
)
```

### メソッドチェーン

```go
record := usage.NewMeteringRecord("product-id", "customer-id", "requests", 100, time.Now()).
    WithCloud("aws").
    WithScheduledAt(time.Now().Add(1 * time.Hour))
```

### ビルダーパターン

```go
record := usage.NewMeteringRecordBuilder().
    ProductID("product-id").
    CustomerID("customer-id").
    Dimension("api-requests", 100).
    StartTime(time.Now()).
    Cloud("aws").
    Build()
```

## API層の構造

このSDKは3つのレベルのAPIを提供します：

1. **高レベルAPI（推奨）**: `client.V1().CreateMeteringUsageRecords()`
2. **中レベルAPI**: `client.CreateUsageRecords()`
3. **低レベルAPI**: 生成されたクライアント直接アクセス

## 複数レコードの一括送信

```go
records := []usage.MeteringRecord{
    usage.NewMeteringRecordWithCloud("aws", "prod-1", "cust-1", "api-calls", 100, time.Now()),
    usage.NewMeteringRecordWithCloud("aws", "prod-1", "cust-2", "api-calls", 150, time.Now()),
    usage.NewMeteringRecordBuilder().
        ProductID("prod-2").
        CustomerID("cust-1").
        Dimension("compute-hours", 5).
        StartTime(time.Now()).
        Cloud("aws").
        Build(),
}

response, err := client.V1().CreateMeteringUsageRecords(ctx, records)
```

## エラーハンドリング

```go
response, err := client.V1().CreateMeteringUsageRecords(ctx, records)
if err != nil {
    log.Printf("API Error: %v", err)
    return
}

switch response.StatusCode() {
case 200:
    if response.JSON200 != nil && response.JSON200.Message != nil {
        log.Printf("Success: %s", *response.JSON200.Message)
    }
case 409:
    if response.JSON409 != nil && response.JSON409.Error != nil {
        log.Printf("Conflict: %s", *response.JSON409.Error)
    }
default:
    log.Printf("Unexpected status: %d", response.StatusCode())
}
```

## 開発

### テスト実行

```bash
go test -v ./...
```

### ビルド

```bash
go build ./...
```

## ライセンス

[ライセンス情報をここに記載]

## 貢献

プルリクエストやイシューの報告を歓迎します。

## サポート

- [API仕様](./api.yaml)
- [開発者ガイド](./CLAUDE.md)