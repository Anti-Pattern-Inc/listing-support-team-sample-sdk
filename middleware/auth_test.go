package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestNewAuthConfigFromEnv(t *testing.T) {
	// Save original env vars
	originalAPIKey := os.Getenv("AWS_MARKETPLACE_API_KEY")
	originalSecretKey := os.Getenv("AWS_MARKETPLACE_SECRET_KEY")
	defer func() {
		os.Setenv("AWS_MARKETPLACE_API_KEY", originalAPIKey)
		os.Setenv("AWS_MARKETPLACE_SECRET_KEY", originalSecretKey)
	}()

	// Test with both keys - should use signature auth
	os.Setenv("AWS_MARKETPLACE_API_KEY", "test-api-key")
	os.Setenv("AWS_MARKETPLACE_SECRET_KEY", "test-secret-key")
	
	config := NewAuthConfigFromEnv()
	if config == nil {
		t.Fatal("expected config, got nil")
	}
	if config.Type != AuthTypeSignature {
		t.Errorf("expected AuthTypeSignature, got %v", config.Type)
	}
	if config.APIKey != "test-api-key" {
		t.Errorf("expected test-api-key, got %v", config.APIKey)
	}
	if config.SecretKey != "test-secret-key" {
		t.Errorf("expected test-secret-key, got %v", config.SecretKey)
	}

	// Test with only API key - should use bearer auth
	os.Setenv("AWS_MARKETPLACE_SECRET_KEY", "")
	
	config = NewAuthConfigFromEnv()
	if config == nil {
		t.Fatal("expected config, got nil")
	}
	if config.Type != AuthTypeBearer {
		t.Errorf("expected AuthTypeBearer, got %v", config.Type)
	}

	// Test with no keys - should return nil
	os.Setenv("AWS_MARKETPLACE_API_KEY", "")
	
	config = NewAuthConfigFromEnv()
	if config != nil {
		t.Errorf("expected nil, got %v", config)
	}
}

func TestValidateAuthConfig(t *testing.T) {
	// Test nil config
	err := ValidateAuthConfig(nil)
	if err == nil {
		t.Error("expected error for nil config")
	}

	// Test valid Bearer config
	config := &AuthConfig{
		Type:   AuthTypeBearer,
		APIKey: "test-key",
	}
	err = ValidateAuthConfig(config)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test invalid Bearer config (no API key)
	config = &AuthConfig{
		Type: AuthTypeBearer,
	}
	err = ValidateAuthConfig(config)
	if err == nil {
		t.Error("expected error for Bearer config without API key")
	}

	// Test valid Signature config
	config = &AuthConfig{
		Type:      AuthTypeSignature,
		APIKey:    "test-key",
		SecretKey: "test-secret",
	}
	err = ValidateAuthConfig(config)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test invalid Signature config (no secret key)
	config = &AuthConfig{
		Type:   AuthTypeSignature,
		APIKey: "test-key",
	}
	err = ValidateAuthConfig(config)
	if err == nil {
		t.Error("expected error for Signature config without secret key")
	}

	// Test valid Custom config
	config = &AuthConfig{
		Type: AuthTypeCustom,
		CustomAuth: func(ctx context.Context, req *http.Request) error {
			return nil
		},
	}
	err = ValidateAuthConfig(config)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// Test invalid Custom config (no handler)
	config = &AuthConfig{
		Type: AuthTypeCustom,
	}
	err = ValidateAuthConfig(config)
	if err == nil {
		t.Error("expected error for Custom config without handler")
	}
}

func TestBearerAuth(t *testing.T) {
	config := &AuthConfig{
		Type:   AuthTypeBearer,
		APIKey: "test-bearer-token",
	}

	requestEditor, err := config.CreateRequestEditor()
	if err != nil {
		t.Fatalf("failed to create request editor: %v", err)
	}

	req, _ := http.NewRequest("GET", "https://example.com/test", nil)
	err = requestEditor(context.Background(), req)
	if err != nil {
		t.Fatalf("request editor failed: %v", err)
	}

	authHeader := req.Header.Get("Authorization")
	expectedAuth := "Bearer test-bearer-token"
	if authHeader != expectedAuth {
		t.Errorf("expected %q, got %q", expectedAuth, authHeader)
	}

	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected application/json, got %q", contentType)
	}
}

func TestSignatureAuth(t *testing.T) {
	config := &AuthConfig{
		Type:      AuthTypeSignature,
		APIKey:    "test-api-key",
		SecretKey: "test-secret-key",
	}

	requestEditor, err := config.CreateRequestEditor()
	if err != nil {
		t.Fatalf("failed to create request editor: %v", err)
	}

	req, _ := http.NewRequest("GET", "https://example.com/test", nil)
	err = requestEditor(context.Background(), req)
	if err != nil {
		t.Fatalf("request editor failed: %v", err)
	}

	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		t.Error("expected Authorization header to be set")
	}

	// Check that it starts with the expected literal
	if !strings.HasPrefix(authHeader, "AWSMARKETPLACESIGV1") {
		t.Errorf("expected Authorization header to start with AWSMARKETPLACESIGV1, got %q", authHeader)
	}

	// Check that it contains expected components
	if !strings.Contains(authHeader, "Sig=") {
		t.Error("expected Authorization header to contain Sig=")
	}
	if !strings.Contains(authHeader, "APIKey=test-api-key") {
		t.Error("expected Authorization header to contain APIKey=test-api-key")
	}
	if !strings.Contains(authHeader, "Timestamp=") {
		t.Error("expected Authorization header to contain Timestamp=")
	}
}

func TestCustomAuth(t *testing.T) {
	customHandler := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("X-Custom-Auth", "custom-value")
		return nil
	}

	config := &AuthConfig{
		Type:       AuthTypeCustom,
		CustomAuth: customHandler,
	}

	requestEditor, err := config.CreateRequestEditor()
	if err != nil {
		t.Fatalf("failed to create request editor: %v", err)
	}

	req, _ := http.NewRequest("GET", "https://example.com/test", nil)
	err = requestEditor(context.Background(), req)
	if err != nil {
		t.Fatalf("request editor failed: %v", err)
	}

	customHeader := req.Header.Get("X-Custom-Auth")
	if customHeader != "custom-value" {
		t.Errorf("expected custom-value, got %q", customHeader)
	}
}