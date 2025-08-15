# CLAUDE.md

このファイルは、このリポジトリでコードを操作する際にClaude Code (claude.ai/code) にガイダンスを提供します。

## プロジェクト概要

このプロジェクトは、AWS Marketplace SaaS使用量レコードAPI用のGo SDKです。OpenAPIから自動生成されたクライアントと、SaaSusスタイルの使いやすいラッパーAPIを提供します。

### ファイル構造

- `openapi.yaml`: OpenAPI 3.0.3仕様書
- `generate.sh`: oapi-codegenでクライアントコード生成スクリプト
- `go.mod`: モジュール定義（`github.com/Anti-Pattern-Inc/listing-support-team-sample-sdk`）
- `generated/clientapi/`: 自動生成されたAPIクライアント
  - `client.gen.go`: HTTPクライアントの実装
  - `types.gen.go`: API型定義
- `middleware/`: HTTPミドルウェア
  - `authenticate.go`: 認証ミドルウェア
  - `example.go`: ミドルウェア使用例
- `modules/clientapi/`: 手動作成の高レベルラッパー
  - `auth.go`: 認証クライアント
  - `client.go`: マーケットプレイスクライアント
  - `env.go`: 環境変数認証

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
# openapi.yamlからクライアントコードを再生成
./generate.sh
```

## 重要な設計原則

### アーキテクチャ構造

1. **generated/clientapi** (基盤レイヤー)
   - OpenAPI仕様からoapi-codegenで自動生成
   - 低レベルHTTPクライアント機能
   - 手動編集禁止（DO NOT EDITコメント付き）

2. **modules/clientapi** (アプリケーションレイヤー)
   - generated/clientapiを利用したラッパー
   - SaaSusスタイルのシンプルなAPI提供
   - 認証処理の自動化
   - ユーザーフレンドリーなインターフェース

### 設計パターン

#### SaaSusスタイルAPI設計
TypeScript SaaSus SDKのパターンを参考に、以下の特徴を持つ：

```go
// TypeScript SaaSus SDK風の使い方
authClient, err := marketplace.NewAuthClient()
token, err := authClient.GetAuthToken(ctx, clientID, apiKey)

marketplaceClient, err := marketplace.NewMarketplaceClient(token)
_, response, err := marketplaceClient.CreateUsageRecords(ctx, records)
```

#### 利便性の優先
- 認証ボイラープレートの排除
- ワンライナー初期化サポート
- 責務分離（認証 vs ビジネスロジック）

### パッケージ構造

#### modules/clientapi パッケージ
- `AuthClient`: 認証専用クライアント
  - `NewAuthClient()`: 認証クライアント作成
  - `GetAuthToken(ctx, clientID, apiKey)`: トークン取得
- `MarketplaceClient`: マーケットプレイス操作クライアント
  - `NewMarketplaceClientFromEnv(ctx)`: 環境変数から認証（推奨）
  - `NewMarketplaceClient(token)`: トークンからクライアント作成
  - `NewMarketplaceClientWithAuth(ctx, clientID, apiKey)`: 認証込み作成
  - `Client`: 全APIメソッドへの直接アクセス

#### generated/clientapi パッケージ
- `Client`: 低レベルHTTPクライアント
- `ClientWithResponses`: レスポンス付きクライアント
- `UsageRecord`: 使用量レコード型
- `Dimension`: ディメンション型
- `UsageAllocation`: 使用量アロケーション型

### 生成コード管理

#### 自動生成ファイル
- `generated/clientapi/client.gen.go`
- `generated/clientapi/types.gen.go`
- **重要**: これらのファイルは手動編集禁止
- 再生成時は `./generate.sh` を実行

#### 手動作成ファイル
- `modules/clientapi/auth.go`
- `modules/clientapi/client.go`
- 設計方針に従って自由に編集可能

## 使用例パターン

### 推奨パターン（SaaSusスタイル）

```go
// パターン1: 環境変数認証（最推奨）
client, err := marketplace.NewMarketplaceClientFromEnv(ctx)
response, err := client.Client.CreateUsageRecordsWithResponse(ctx, records)

// パターン2: 直接認証
client, err := marketplace.NewMarketplaceClientWithAuth(ctx, clientID, apiKey)
response, err := client.Client.CreateUsageRecordsWithResponse(ctx, records)

// パターン3: 段階的認証（高度な制御）
authClient, err := marketplace.NewAuthClient()
token, err := authClient.GetAuthToken(ctx, clientID, apiKey)
marketplaceClient, err := marketplace.NewMarketplaceClient(token)
response, err := marketplaceClient.Client.CreateUsageRecordsWithResponse(ctx, records)
```

### 低レベルパターン（直接generated使用）

```go
// 認証トークン取得
client, err := clientapi.NewClientWithResponses(serverURL)
authRequest := clientapi.IssueAuthTokenJSONRequestBody{
    ClientId: clientID,
    ApiKey:   apiKey,
}
response, err := client.IssueAuthTokenWithResponse(ctx, authRequest)
token := *response.JSON200.AccessToken

// 認証済みクライアント作成
authFunc := func(ctx context.Context, req *http.Request) error {
    req.Header.Set("Authorization", "Bearer "+token)
    return nil
}
authenticatedClient, err := clientapi.NewClientWithResponses(
    serverURL,
    clientapi.WithRequestEditorFn(authFunc),
)

// 使用量記録作成
records := clientapi.CreateUsageRecordsJSONRequestBody{...}
response, err := authenticatedClient.CreateUsageRecordsWithResponse(ctx, records)
```

## API仕様

### 認証フロー
1. **トークン取得**: `POST /v1/auth/token` でClientID + API_KEYからJWTトークン取得
2. **API呼び出し**: 取得したトークンを `Authorization: Bearer {jwt-token}` で送信
3. **トークン有効期限**: 1時間

### エンドポイント
- `POST /v1/auth/token`: 認証トークン発行
- `POST /v1/usage-records`: 使用量記録作成
- `PUT /v1/usage-records`: 使用量記録更新
- `OPTIONS /v1/usage-records`: CORS対応

### 主要な型

#### UsageRecord
```go
type UsageRecord struct {
    ProductId          string    `json:"product_id"`
    CustomerIdentifier string    `json:"CustomerIdentifier"`
    Dimension          Dimension `json:"dimension"`
    StartTime          time.Time `json:"start_time"`
}
```

#### Dimension
```go
type Dimension struct {
    Name             string             `json:"name"`
    Quantity         float32            `json:"quantity"`
    DryRun          *bool              `json:"dry_run,omitempty"`
    UsageAllocations *[]UsageAllocation `json:"usage_allocations,omitempty"`
}
```

## 開発ガイドライン

### コーディング規約
1. **ネストは3階層まで**
2. **変数名、関数名は意味のある名前**
3. **不要なファイルは作成しない**
4. **コメントは英語で記述**

### エラーハンドリング
- すべてのAPIコールでエラーチェック実装
- ステータスコード別の適切な処理
- ユーザーフレンドリーなエラーメッセージ

### テスト
- 単体テストの充実
- 認証フローのテスト
- エラーケースの網羅

## 保守性の考慮

### 関心の分離
- **generated**: OpenAPI仕様からの純粋なコード生成
- **modules**: ビジネスロジックと認証の処理

### 拡張性
- OpenAPI仕様更新時はgeneratedのみ再生成
- modulesのカスタムロジックは影響を受けない
- 新しい認証方式の追加が容易

### 使いやすさ
- エンドユーザーはmodulesのAPIを使用
- 認証やエンドポイント設定が自動化
- TypeScript SaaSus SDKと同様の開発体験