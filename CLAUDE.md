# CLAUDE.md

このファイルは、このリポジトリでコードを操作する際にClaude Code (claude.ai/code) にガイダンスを提供します。

## プロジェクト概要

このプロジェクトは、AWS Marketplace SaaS使用量レコードAPI用のGo SDKです。

### ファイル構造

- `client.go`: oapi-codegenで自動生成されたAPIクライアント
- `sdk.go`: ユーザーフレンドリーなSDKラッパー
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
- パッケージ名: `sdk`
- `sdk.Config`: SDK設定（BaseURL, APIKey, HTTPClient）
- `sdk.SDK`: メインSDK構造体
- `sdk.NewSDK()`: SDKインスタンス作成
- Helper関数: `NewUsageRecord()`, `NewUsageAllocation()`, `NewTag()`

### 生成コード
- `client.go`はoapi-codegenで自動生成
- 生成後は必ずパッケージ名を`sdk`に変更
- 手動編集は避け、再生成時はパッケージ名の修正を忘れずに

### 使用例
```go
import "github.com/nishimurashinya/openapi-go-sdk"

config := &sdk.Config{
    BaseURL: "https://api-id.execute-api.ap-northeast-1.amazonaws.com/Prod",
    APIKey:  "your-api-key",
}

client, err := sdk.NewSDK(config)
if err != nil {
    log.Fatal(err)
}

record := sdk.NewUsageRecord("product-id", "customer-id", "dimension", 100, time.Now())
response, err := client.CreateUsageRecords(ctx, []sdk.UsageRecord{record})
```

## API仕様
- OpenAPI 3.0.3仕様（api.yaml）
- AWS Marketplace SaaS使用量レコードAPI
- Bearer認証使用
- DynamoDBとSQSを使用したサーバーレス構成