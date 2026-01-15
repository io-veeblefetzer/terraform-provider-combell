// Copyright (c) Veeblefetzer
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name      string
		apiKey    string
		apiSecret string
		baseURL   string
		wantURL   string
	}{
		{
			name:      "default base URL",
			apiKey:    "test-key",
			apiSecret: "test-secret",
			baseURL:   "",
			wantURL:   DefaultBaseURL,
		},
		{
			name:      "custom base URL",
			apiKey:    "test-key",
			apiSecret: "test-secret",
			baseURL:   "https://custom.api.com/v2",
			wantURL:   "https://custom.api.com/v2",
		},
		{
			name:      "custom base URL with trailing slash",
			apiKey:    "test-key",
			apiSecret: "test-secret",
			baseURL:   "https://custom.api.com/v2/",
			wantURL:   "https://custom.api.com/v2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.apiKey, tt.apiSecret, tt.baseURL)

			if client.BaseURL != tt.wantURL {
				t.Errorf("NewClient() BaseURL = %v, want %v", client.BaseURL, tt.wantURL)
			}

			if client.APIKey != tt.apiKey {
				t.Errorf("NewClient() APIKey = %v, want %v", client.APIKey, tt.apiKey)
			}

			if client.APISecret != tt.apiSecret {
				t.Errorf("NewClient() APISecret = %v, want %v", client.APISecret, tt.apiSecret)
			}

			if client.HTTPClient == nil {
				t.Error("NewClient() HTTPClient is nil")
			}
		})
	}
}

func TestGenerateHMAC(t *testing.T) {
	client := NewClient("my-api-key", "my-api-secret", "")

	// Test HMAC generation with known values
	// The exact signature will depend on the implementation details
	signature := client.generateHMAC("GET", "/v2/dns/example.com/records", 1609459200, "unique-nonce-123", nil)

	if signature == "" {
		t.Error("generateHMAC() returned empty signature")
	}

	// Test that same inputs produce same signature
	signature2 := client.generateHMAC("GET", "/v2/dns/example.com/records", 1609459200, "unique-nonce-123", nil)

	if signature != signature2 {
		t.Error("generateHMAC() not deterministic - same inputs produced different signatures")
	}

	// Test that different nonce produces different signature
	signature3 := client.generateHMAC("GET", "/v2/dns/example.com/records", 1609459200, "different-nonce", nil)

	if signature == signature3 {
		t.Error("generateHMAC() same signature with different nonce")
	}

	// Test with body content
	body := []byte(`{"type":"A","content":"192.168.1.1"}`)
	signatureWithBody := client.generateHMAC("POST", "/v2/dns/example.com/records", 1609459200, "unique-nonce-123", body)

	if signatureWithBody == "" {
		t.Error("generateHMAC() with body returned empty signature")
	}

	if signatureWithBody == signature {
		t.Error("generateHMAC() same signature with and without body")
	}
}

func TestEncodePathForHMAC(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "simple path",
			path:     "/v2/dns/example.com/records",
			expected: "/v2/dns/example.com/records",
		},
		{
			name:     "path with uppercase",
			path:     "/v2/DNS/Example.Com/Records",
			expected: "/v2/dns/example.com/records",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := encodePathForHMAC(tt.path)
			if result != tt.expected {
				t.Errorf("encodePathForHMAC(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

func TestUppercasePercentEncoding(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no encoding",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "lowercase percent encoding",
			input:    "%2f%3a%40",
			expected: "%2F%3A%40",
		},
		{
			name:     "mixed case",
			input:    "%2F%3a%40",
			expected: "%2F%3A%40",
		},
		{
			name:     "with regular text",
			input:    "hello%2fworld",
			expected: "hello%2Fworld",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uppercasePercentEncoding(tt.input)
			if result != tt.expected {
				t.Errorf("uppercasePercentEncoding(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBuildAuthorizationHeader(t *testing.T) {
	client := NewClient("my-api-key", "my-api-secret", "")

	header := client.buildAuthorizationHeader("signature123", "nonce456", 1609459200)

	expected := "hmac my-api-key:signature123:nonce456:1609459200"
	if header != expected {
		t.Errorf("buildAuthorizationHeader() = %q, want %q", header, expected)
	}
}

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *APIError
		contains []string
	}{
		{
			name: "simple error",
			err: &APIError{
				StatusCode: 404,
				Message:    "not found",
			},
			contains: []string{"404", "not found"},
		},
		{
			name: "validation error",
			err: &APIError{
				StatusCode: 400,
				ValidationErrors: []ValidationError{
					{ErrorCode: "invalid_field", ErrorText: "field is required"},
				},
			},
			contains: []string{"400", "invalid_field", "field is required"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errMsg := tt.err.Error()
			for _, substr := range tt.contains {
				if !containsString(errMsg, substr) {
					t.Errorf("APIError.Error() = %q, want it to contain %q", errMsg, substr)
				}
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStringHelper(s, substr))
}

func containsStringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestClient_Get(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Verify authorization header exists
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			t.Error("Authorization header is missing")
		}

		// Verify it starts with "hmac "
		if len(authHeader) < 5 || authHeader[:5] != "hmac " {
			t.Errorf("Authorization header should start with 'hmac ', got %s", authHeader)
		}

		// Return a successful response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	client := NewClient("test-key", "test-secret", server.URL)

	_, body, err := client.Get(context.Background(), "/test")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if len(body) == 0 {
		t.Error("Get() returned empty body")
	}
}

func TestClient_Post(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Verify Content-Type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Content-Type = %s, want application/json", contentType)
		}

		w.Header().Set("Location", "/created/resource")
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := NewClient("test-key", "test-secret", server.URL)

	resp, _, err := client.Post(context.Background(), "/test", map[string]string{"key": "value"})
	if err != nil {
		t.Fatalf("Post() error = %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Post() status = %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	location := GetLocationHeader(resp)
	if location != "/created/resource" {
		t.Errorf("Location header = %s, want /created/resource", location)
	}
}

func TestClient_ErrorHandling(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		responseBody string
		expectError  bool
	}{
		{
			name:         "success",
			statusCode:   http.StatusOK,
			responseBody: `{"status":"ok"}`,
			expectError:  false,
		},
		{
			name:         "bad request",
			statusCode:   http.StatusBadRequest,
			responseBody: `{"validation_errors":[{"error_code":"invalid","error_text":"invalid field"}]}`,
			expectError:  true,
		},
		{
			name:         "not found",
			statusCode:   http.StatusNotFound,
			responseBody: `{"error":"not found"}`,
			expectError:  true,
		},
		{
			name:         "unauthorized",
			statusCode:   http.StatusUnauthorized,
			responseBody: `{"error":"unauthorized"}`,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client := NewClient("test-key", "test-secret", server.URL)
			_, _, err := client.Get(context.Background(), "/test")

			if tt.expectError && err == nil {
				t.Error("Expected error but got nil")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if tt.expectError && err != nil {
				if apiErr, ok := err.(*APIError); ok {
					if apiErr.StatusCode != tt.statusCode {
						t.Errorf("APIError.StatusCode = %d, want %d", apiErr.StatusCode, tt.statusCode)
					}
				}
			}
		})
	}
}

func TestGetRateLimitInfo(t *testing.T) {
	header := make(http.Header)
	header.Set("X-RateLimit-Limit", "100")
	header.Set("X-RateLimit-Usage", "50")
	header.Set("X-RateLimit-Remaining", "50")
	header.Set("X-RateLimit-Reset", "3600")
	resp := &http.Response{
		Header: header,
	}

	info := GetRateLimitInfo(resp)

	if info.Limit != 100 {
		t.Errorf("Limit = %d, want 100", info.Limit)
	}
	if info.Usage != 50 {
		t.Errorf("Usage = %d, want 50", info.Usage)
	}
	if info.Remaining != 50 {
		t.Errorf("Remaining = %d, want 50", info.Remaining)
	}
	if info.Reset != 3600 {
		t.Errorf("Reset = %d, want 3600", info.Reset)
	}
}

func TestGetRateLimitInfo_NilResponse(t *testing.T) {
	info := GetRateLimitInfo(nil)
	if info != nil {
		t.Error("GetRateLimitInfo(nil) should return nil")
	}
}

func TestGetPagingInfo(t *testing.T) {
	header := make(http.Header)
	header.Set("X-Paging-Skipped", "0")
	header.Set("X-Paging-Take", "25")
	header.Set("X-Paging-TotalResults", "100")
	resp := &http.Response{
		Header: header,
	}

	info := GetPagingInfo(resp)

	if info.Skipped != 0 {
		t.Errorf("Skipped = %d, want 0", info.Skipped)
	}
	if info.Take != 25 {
		t.Errorf("Take = %d, want 25", info.Take)
	}
	if info.TotalResults != 100 {
		t.Errorf("TotalResults = %d, want 100", info.TotalResults)
	}
}

func TestGetLocationHeader(t *testing.T) {
	resp := &http.Response{
		Header: http.Header{
			"Location": []string{"/v2/dns/example.com/records/123"},
		},
	}

	location := GetLocationHeader(resp)
	if location != "/v2/dns/example.com/records/123" {
		t.Errorf("Location = %s, want /v2/dns/example.com/records/123", location)
	}

	// Test nil response
	location = GetLocationHeader(nil)
	if location != "" {
		t.Errorf("GetLocationHeader(nil) = %s, want empty string", location)
	}
}
