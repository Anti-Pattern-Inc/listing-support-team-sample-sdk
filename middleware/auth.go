package middleware

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// AuthType represents the type of authentication
type AuthType string

const (
	// AuthTypeBearer uses Bearer token authentication
	AuthTypeBearer AuthType = "bearer"
	// AuthTypeSignature uses HMAC-SHA256 signature authentication
	AuthTypeSignature AuthType = "signature"
	// AuthTypeCustom allows custom authentication
	AuthTypeCustom AuthType = "custom"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
	// Type specifies the authentication type
	Type AuthType
	// APIKey for Bearer or signature authentication
	APIKey string
	// SecretKey for signature authentication (optional)
	SecretKey string
	// CustomAuth for custom authentication handler
	CustomAuth func(ctx context.Context, req *http.Request) error
}

// NewAuthConfigFromEnv creates AuthConfig from environment variables
func NewAuthConfigFromEnv() *AuthConfig {
	apiKey := os.Getenv("AWS_MARKETPLACE_API_KEY")
	secretKey := os.Getenv("AWS_MARKETPLACE_SECRET_KEY")
	
	// If both keys are available, use signature auth
	if apiKey != "" && secretKey != "" {
		return &AuthConfig{
			Type:      AuthTypeSignature,
			APIKey:    apiKey,
			SecretKey: secretKey,
		}
	}
	
	// Fallback to bearer auth if only API key is available
	if apiKey != "" {
		return &AuthConfig{
			Type:   AuthTypeBearer,
			APIKey: apiKey,
		}
	}
	
	return nil
}

// CreateRequestEditor creates a request editor function based on auth config
func (ac *AuthConfig) CreateRequestEditor() (func(ctx context.Context, req *http.Request) error, error) {
	if ac == nil {
		return nil, fmt.Errorf("auth config is required")
	}

	switch ac.Type {
	case AuthTypeBearer:
		return ac.createBearerAuth(), nil
	case AuthTypeSignature:
		return ac.createSignatureAuth(), nil
	case AuthTypeCustom:
		if ac.CustomAuth == nil {
			return nil, fmt.Errorf("custom auth handler is required for custom auth type")
		}
		return ac.CustomAuth, nil
	default:
		return nil, fmt.Errorf("unsupported auth type: %s", ac.Type)
	}
}

// createBearerAuth creates a Bearer token authentication handler
func (ac *AuthConfig) createBearerAuth() func(ctx context.Context, req *http.Request) error {
	return func(ctx context.Context, req *http.Request) error {
		if ac.APIKey == "" {
			return fmt.Errorf("API key is required for Bearer authentication")
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ac.APIKey))
		req.Header.Set("Content-Type", "application/json")
		return nil
	}
}

// createSignatureAuth creates a signature-based authentication handler
func (ac *AuthConfig) createSignatureAuth() func(ctx context.Context, req *http.Request) error {
	return func(ctx context.Context, req *http.Request) error {
		if ac.APIKey == "" || ac.SecretKey == "" {
			return fmt.Errorf("both API key and secret key are required for signature authentication")
		}
		
		sig, err := ac.generateSignature(req)
		if err != nil {
			return fmt.Errorf("failed to generate signature: %w", err)
		}
		
		req.Header.Set("Authorization", sig)
		req.Header.Set("Content-Type", "application/json")
		return nil
	}
}

// signatureV1 represents the signature generation data
type signatureV1 struct {
	timestamp       string
	secretKey       string
	apiKey          string
	httpMethod      string
	httpURL         string
	httpRequestBody io.ReadCloser
}

// generateSignature generates HMAC-SHA256 signature for request authentication
func (ac *AuthConfig) generateSignature(req *http.Request) (string, error) {
	timestamp := time.Now().UTC().Format("200601021504")
	literal := "AWSMARKETPLACESIGV1"
	
	sig := &signatureV1{
		timestamp:       timestamp,
		secretKey:       ac.SecretKey,
		apiKey:          ac.APIKey,
		httpMethod:      req.Method,
		httpURL:         req.URL.String(),
		httpRequestBody: req.Body,
	}
	
	// Generate HMAC-SHA256 signature
	h := hmac.New(sha256.New, []byte(sig.secretKey))
	h.Write([]byte(sig.timestamp))
	h.Write([]byte(sig.apiKey))
	h.Write([]byte(strings.ToUpper(sig.httpMethod)))
	
	// Extract host and path from URL
	hostURI := strings.Split(sig.httpURL, "//")
	if len(hostURI) < 2 {
		return "", fmt.Errorf("invalid URL format: %s", sig.httpURL)
	}
	h.Write([]byte(hostURI[1]))
	
	// Include request body in signature if present
	if sig.httpRequestBody != nil {
		bodyBytes, err := io.ReadAll(sig.httpRequestBody)
		if err != nil {
			return "", fmt.Errorf("failed to read request body: %w", err)
		}
		h.Write(bodyBytes)
		// Restore the body for the actual request
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}
	
	// Generate final signature
	signature := hex.EncodeToString(h.Sum(nil))
	authHeader := fmt.Sprintf("%s Sig=%s, APIKey=%s, Timestamp=%s", 
		literal, signature, sig.apiKey, sig.timestamp)
	
	return authHeader, nil
}



// ValidateAuthConfig validates the authentication configuration
func ValidateAuthConfig(config *AuthConfig) error {
	if config == nil {
		return fmt.Errorf("auth config cannot be nil")
	}
	
	switch config.Type {
	case AuthTypeBearer:
		if config.APIKey == "" {
			return fmt.Errorf("API key is required for Bearer authentication")
		}
	case AuthTypeSignature:
		if config.APIKey == "" || config.SecretKey == "" {
			return fmt.Errorf("both API key and secret key are required for signature authentication")
		}
	case AuthTypeCustom:
		if config.CustomAuth == nil {
			return fmt.Errorf("custom auth handler is required for custom authentication")
		}
	default:
		return fmt.Errorf("unsupported authentication type: %s", config.Type)
	}
	
	return nil
}