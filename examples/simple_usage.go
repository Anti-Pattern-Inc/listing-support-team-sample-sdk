package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Anti-Pattern-Inc/listing-support-team-sample-sdk/generated/clientapi"
)

func main() {
	ctx := context.Background()

	// 🔥 ここを実際の値に変更してください
	baseURL := "https://your-api-endpoint.com"
	clientID := "your_client_id"
	apiKey := "your_api_key"

	fmt.Println("🚀 API クライアント開始")

	// Step 1: 認証
	fmt.Println("🔐 認証中...")
	token, err := getAuthToken(ctx, baseURL, clientID, apiKey)
	if err != nil {
		log.Fatal("❌ 認証失敗:", err)
	}
	fmt.Println("✅ 認証成功")

	// Step 2: 認証済みクライアント作成
	client, err := createAuthenticatedClient(baseURL, token)
	if err != nil {
		log.Fatal("❌ クライアント作成失敗:", err)
	}

	// Step 3: 使用量記録
	fmt.Println("📊 使用量記録中...")
	err = createUsageRecord(ctx, client)
	if err != nil {
		log.Fatal("❌ 使用量記録失敗:", err)
	}
	fmt.Println("✅ 完了!")
}

func getAuthToken(ctx context.Context, baseURL, clientID, apiKey string) (string, error) {
	client, err := clientapi.NewClientWithResponses(baseURL)
	if err != nil {
		return "", err
	}

	authRequest := clientapi.IssueAuthTokenJSONRequestBody{
		ClientId: clientID,
		ApiKey:   apiKey,
	}

	response, err := client.IssueAuthTokenWithResponse(ctx, authRequest)
	if err != nil {
		return "", err
	}

	switch response.StatusCode() {
	case 200:
		if response.JSON200 != nil && response.JSON200.AccessToken != nil {
			return *response.JSON200.AccessToken, nil
		}
		return "", fmt.Errorf("アクセストークンが見つかりません")
	case 401:
		return "", fmt.Errorf("認証情報が無効: %s", response.JSON401.Error)
	case 500:
		return "", fmt.Errorf("サーバーエラー: %s", response.JSON500.Error)
	default:
		return "", fmt.Errorf("予期しないエラー: ステータス %d", response.StatusCode())
	}
}

func createAuthenticatedClient(baseURL, token string) (*clientapi.ClientWithResponses, error) {
	authFunc := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	}

	return clientapi.NewClientWithResponses(
		baseURL,
		clientapi.WithRequestEditorFn(authFunc),
	)
}

func createUsageRecord(ctx context.Context, client *clientapi.ClientWithResponses) error {
	records := clientapi.CreateUsageRecordsJSONRequestBody{
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

	response, err := client.CreateUsageRecordsWithResponse(ctx, records)
	if err != nil {
		return err
	}

	switch response.StatusCode() {
	case 200:
		if response.JSON200 != nil && response.JSON200.Message != nil {
			fmt.Printf("📈 記録成功: %s\n", *response.JSON200.Message)
		} else {
			fmt.Println("📈 記録成功")
		}
		return nil
	case 401:
		return fmt.Errorf("認証エラー: %s", response.JSON401.Error)
	case 400:
		return fmt.Errorf("リクエストエラー: %s", response.JSON400.Error)
	case 500:
		return fmt.Errorf("サーバーエラー: %s", response.JSON500.Error)
	default:
		return fmt.Errorf("予期しないエラー: ステータス %d", response.StatusCode())
	}
}
