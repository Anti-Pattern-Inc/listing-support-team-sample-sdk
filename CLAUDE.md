# CLAUDE.md

このファイルは、このリポジトリでコードを操作する際にClaude Code (claude.ai/code) にガイダンスを提供します。

## プロジェクト概要

このプロジェクトは、AWS Marketplace SaaS使用量レコードAPI用のGo SDKです。

### ファイル構造

- `client.go`: oapi-codegenで自動生成されたAPIクライアント
- `sdk.go`: ユーザーフレンドリーなSDKラッパー
- `wrapper.go`: Node.js風の高レベルAPIラッパー
- `middleware/`: 認証とミドルウェア関連
  - `auth.go`: 認証処理（Bearer, HMAC-SHA256, Custom）
  - `auth_test.go`: 認証機能のテスト
- `wrapper_test.go`: ラッパー機能のテスト
- `example_test.go`: 使用例とサンプルコード
- `go.mod`: モジュール定義（`github.com/nishimurashinya/openapi-go-sdk`）
- `api.yaml`: OpenAPI 3.0.3仕様書

## 開発コマンド

### ビルドとテスト
```bash
# SDKのテスト
go test ./...

# 依存関係の更新
go mod tidy

# SDKのビルド確認
go build ./...
```

### コード生成
```bash
# api.yamlからクライアントコードを再生成する場合
# oapi-codegenを使用してclient.goを生成
# 生成後にパッケージ名をsdkに変更する必要あり
```

## 重要な設計原則

### パッケージ構造
- **メインパッケージ**: `sdk`
  - `sdk.Config`: SDK設定（BaseURL, Auth, HTTPClient）
  - `sdk.SDK`: メインSDK構造体
  - `sdk.MeteringRecord`: Node.js風の簡素化されたレコード構造
  - `sdk.MeteringUsageRecords`: 高レベルAPIインターフェース
  - Helper関数: `NewMeteringRecord()`, `NewMeteringRecordBuilder()`, etc.
- **認証パッケージ**: `middleware`
  - `middleware.AuthConfig`: 認証設定（Type, APIKey, SecretKey, CustomAuth）
  - `middleware.AuthType`: 認証方式の定数

### 認証方式
- `middleware.AuthTypeBearer`: シンプルなBearer認証
- `middleware.AuthTypeSignature`: HMAC-SHA256署名認証（高セキュリティ）
- `middleware.AuthTypeCustom`: カスタム認証ハンドラー

### API層の構造
1. **高レベルAPI**: `client.V1().CreateMeteringUsageRecords()` (推奨)
2. **中レベルAPI**: `client.CreateUsageRecords()`
3. **低レベルAPI**: 生成されたクライアント直接アクセス

### 便利なコンストラクタ
- `NewSDKWithBearer()`: Bearer認証用
- `NewSDKWithSignature()`: 署名認証用
- `NewSDKWithCustomAuth()`: カスタム認証用
- `NewSDKFromEnv()`: 環境変数から自動設定
- `NewMeteringRecord()`: シンプルなレコード作成
- `NewMeteringRecordBuilder()`: ビルダーパターン

### 生成コード
- `client.go`はoapi-codegenで自動生成
- 生成後は必ずパッケージ名を`sdk`に変更
- 手動編集は避け、再生成時はパッケージ名の修正を忘れずに

### 使用例

#### 1. Node.js風シンプルAPI（推奨）
```go
import "github.com/nishimurashinya/openapi-go-sdk"

// SDK作成
client, err := sdk.NewSDKWithBearer(
    "https://api-id.execute-api.ap-northeast-1.amazonaws.com/Prod",
    "your-api-key",
)

// Node.js風の使用量レコード作成
records := []sdk.MeteringRecord{
    {
        Cloud:      "aws",
        ProductID:  "1test1a2b3c4d5e6f7g8h9ijk",
        CustomerID: "ij3sXMkN3or",
        Dimension:  sdk.MeteringDimension{Name: "metering-test-sku-from-vendor-123", Quantity: 1},
        StartTime:  "2019-07-15T15:00:00.000000+00:00",
        ScheduledAt: "2019-07-15T15:00:00.000000+00:00",
    },
}

response, err := client.V1().CreateMeteringUsageRecords(ctx, records)
```

#### 2. 便利なヘルパー関数
```go
// 基本的なレコード作成
record := sdk.NewMeteringRecord("product-id", "customer-id", "dimension", 100, time.Now())

// クラウドプラットフォーム指定
record := sdk.NewMeteringRecordWithCloud("aws", "product-id", "customer-id", "dimension", 100, time.Now())

// スケジュール指定
record := sdk.NewMeteringRecordWithSchedule("product-id", "customer-id", "dimension", 100, startTime, scheduledAt)

// メソッドチェーン
record := sdk.NewMeteringRecord("product-id", "customer-id", "dimension", 100, time.Now()).
    WithCloud("aws").
    WithScheduledAt(scheduledAt)

// ビルダーパターン
record := sdk.NewMeteringRecordBuilder().
    ProductID("product-id").
    CustomerID("customer-id").
    Dimension("dimension", 100).
    StartTime(time.Now()).
    Cloud("aws").
    Build()
```

#### 3. Bearer認証（従来方式）
```go
import (
    "github.com/nishimurashinya/openapi-go-sdk"
    "github.com/nishimurashinya/openapi-go-sdk/middleware"
)

// 簡単な方法
client, err := sdk.NewSDKWithBearer(
    "https://api-id.execute-api.ap-northeast-1.amazonaws.com/Prod",
    "your-api-key",
)

// または設定を使用
config := &sdk.Config{
    BaseURL: "https://api-id.execute-api.ap-northeast-1.amazonaws.com/Prod",
    Auth: &middleware.AuthConfig{
        Type:   middleware.AuthTypeBearer,
        APIKey: "your-api-key",
    },
}
client, err := sdk.NewSDK(config)
```

#### 4. HMAC-SHA256署名認証（高セキュリティ）
```go
// 直接指定
client, err := sdk.NewSDKWithSignature(
    "https://api-id.execute-api.ap-northeast-1.amazonaws.com/Prod",
    "your-api-key",
    "your-secret-key",
)

// 環境変数から自動取得
// AWS_MARKETPLACE_API_KEY および AWS_MARKETPLACE_SECRET_KEY
client, err := sdk.NewSDKFromEnv(
    "https://api-id.execute-api.ap-northeast-1.amazonaws.com/Prod",
)
```

#### 5. カスタム認証
```go
customAuth := func(ctx context.Context, req *http.Request) error {
    req.Header.Set("X-Custom-Token", "your-custom-token")
    return nil
}

client, err := sdk.NewSDKWithCustomAuth(
    "https://api-id.execute-api.ap-northeast-1.amazonaws.com/Prod",
    customAuth,
)
```

#### 6. 低レベルAPI（直接アクセス）
```go
record := sdk.NewUsageRecord("product-id", "customer-id", "dimension", 100, time.Now())
response, err := client.CreateUsageRecords(ctx, []sdk.UsageRecord{record})
```

## API仕様
- OpenAPI 3.0.3仕様（api.yaml）
- AWS Marketplace SaaS使用量レコードAPI
- Bearer認証使用
- DynamoDBとSQSを使用したサーバーレス構成