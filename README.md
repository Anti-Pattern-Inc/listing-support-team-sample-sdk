# AWS Marketplace SaaS Usage Records SDK for Go

AWS Marketplace SaaSサービス向けの使用量レコードAPI用Go SDKです。

## インストール

```bash
go get github.com/Anti-Pattern-Inc/listing-support-team-sample-sdk
```

## 環境変数設定

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
