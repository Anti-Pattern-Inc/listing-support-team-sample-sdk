# AWS Marketplace SaaS Usage Records SDK for Go

AWS Marketplace SaaSサービス向けの使用量レコードAPI用Go SDKです。

## インストール

```bash
go get github.com/Anti-Pattern-Inc/listing-support-team-sample-sdk
```

## 環境変数設定

> **注意**: 環境変数の値は、Sandbox for AWS Marketplace Sellerの`metering-app-api-keys`テーブルから取得してください。
> - `CLIENT_ID`: テーブルの`customer_id`列を参照
> - `API_KEY`: テーブルの`api_key`列を参照

### 方法1: .envファイルを使用（推奨）

1. `.env.example` をコピーして `.env` ファイルを作成：
```bash
cp .env.example .env
```

2. `.env` ファイルを編集して実際の値を設定：
```bash
CLIENT_ID=your_client_id
API_KEY=your_api_key
API_URL=https://your-api-gateway-url.com/Prod
```

3. プログラムが自動的に `.env` ファイルを読み込みます（godotenv使用）

### 方法2: 環境変数を直接設定

```bash
export CLIENT_ID="your_client_id"
export API_KEY="your_api_key"
export API_URL="https://your-api-gateway-url.com/Prod"
```

## 使い方

`example/main.go` を参照してください。

```bash
go run example/main.go
```

## 開発

### クライアント再生成

```bash
./generate.sh
```

### ビルド

```bash
go build ./...
```
