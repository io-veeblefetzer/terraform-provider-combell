// Copyright (c) Veeblefetzer
// SPDX-License-Identifier: MPL-2.0

package client

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
)

const (
	// DefaultBaseURL is the default Combell API base URL
	DefaultBaseURL = "https://api.combell.com/v2"

	// DefaultTimeout is the default HTTP client timeout
	DefaultTimeout = 30 * time.Second

	// MaxRetries is the maximum number of retries for rate-limited or server error requests
	MaxRetries = 3

	// RetryWaitMin is the minimum time to wait between retries
	RetryWaitMin = 1 * time.Second

	// RetryWaitMax is the maximum time to wait between retries
	RetryWaitMax = 30 * time.Second
)

// Client is the Combell API client
type Client struct {
	// BaseURL is the API base URL
	BaseURL string

	// APIKey is the API key for authentication
	APIKey string

	// APISecret is the API secret for HMAC signature
	APISecret string

	// HTTPClient is the underlying HTTP client
	HTTPClient *http.Client

	// UserAgent is the User-Agent header value
	UserAgent string
}

// RateLimitInfo contains rate limit information from API response headers
type RateLimitInfo struct {
	Limit     int
	Usage     int
	Remaining int
	Reset     int
}

// APIError represents an error response from the Combell API
type APIError struct {
	StatusCode       int
	Message          string
	ValidationErrors []ValidationError
}

// ValidationError represents a validation error from the API
type ValidationError struct {
	ErrorCode string `json:"error_code"`
	ErrorText string `json:"error_text"`
}

// Error returns the error message
func (e *APIError) Error() string {
	if len(e.ValidationErrors) > 0 {
		var msgs []string
		for _, ve := range e.ValidationErrors {
			msgs = append(msgs, fmt.Sprintf("%s: %s", ve.ErrorCode, ve.ErrorText))
		}
		return fmt.Sprintf("API error (status %d): %s", e.StatusCode, strings.Join(msgs, "; "))
	}
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Message)
}

// NewClient creates a new Combell API client
func NewClient(apiKey, apiSecret, baseURL string) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	// Ensure base URL doesn't have trailing slash
	baseURL = strings.TrimSuffix(baseURL, "/")

	return &Client{
		BaseURL:   baseURL,
		APIKey:    apiKey,
		APISecret: apiSecret,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		UserAgent: "terraform-provider-combell",
	}
}

// generateHMAC generates the HMAC signature for a request
// According to Combell API docs:
// 1. Concatenate: apikey + method + path + timestamp + nonce + content_hash
// 2. Hash with SHA-256 using the API secret
// 3. Base64 encode the result
func (c *Client) generateHMAC(method, path string, timestamp int64, nonce string, body []byte) string {
	// Build the input string
	// Path must be URL-encoded and lowercased, starting with /v2
	encodedPath := encodePathForHMAC(path)

	// Content hash: base64 encoded MD5 hash of request body (empty if no body)
	var contentHash string
	if len(body) > 0 {
		md5Hash := md5.Sum(body)
		contentHash = base64.StdEncoding.EncodeToString(md5Hash[:])
	}

	// Concatenate all parts
	input := c.APIKey + strings.ToLower(method) + encodedPath + strconv.FormatInt(timestamp, 10) + nonce + contentHash

	// Create HMAC-SHA256 signature
	h := hmac.New(sha256.New, []byte(c.APISecret))
	h.Write([]byte(input))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature
}

// encodePathForHMAC encodes the path according to Combell API requirements
// The path MUST start with the api version (/v2)
// The hexadecimal codes (percent encoding) MUST be uppercased
func encodePathForHMAC(path string) string {
	// Parse the path to separate path and query string
	u, err := url.Parse(path)
	if err != nil {
		return strings.ToLower(path)
	}

	// Encode path segments
	encodedPath := strings.ToLower(u.EscapedPath())

	// Encode query string if present
	if u.RawQuery != "" {
		// URL encode the query string
		encodedQuery := url.QueryEscape(strings.ToLower("?" + u.RawQuery))
		// Uppercase the percent-encoded characters
		encodedQuery = uppercasePercentEncoding(encodedQuery)
		encodedPath += encodedQuery
	}

	return encodedPath
}

// uppercasePercentEncoding converts percent-encoded characters to uppercase
func uppercasePercentEncoding(s string) string {
	var result strings.Builder
	i := 0
	for i < len(s) {
		if s[i] == '%' && i+2 < len(s) {
			result.WriteByte('%')
			result.WriteString(strings.ToUpper(s[i+1 : i+3]))
			i += 3
		} else {
			result.WriteByte(s[i])
			i++
		}
	}
	return result.String()
}

// buildAuthorizationHeader builds the Authorization header value
// Format: hmac apikey:signature:nonce:timestamp
func (c *Client) buildAuthorizationHeader(signature, nonce string, timestamp int64) string {
	return fmt.Sprintf("hmac %s:%s:%s:%d", c.APIKey, signature, nonce, timestamp)
}

// doRequest performs an HTTP request with HMAC authentication
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}) (*http.Response, []byte, error) {
	var bodyBytes []byte
	var err error

	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	// Build full URL
	fullURL := c.BaseURL + path

	// Generate authentication parameters
	timestamp := time.Now().Unix()
	nonce, err := uuid.GenerateUUID()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Generate HMAC signature (path for signature should not include base URL)
	signature := c.generateHMAC(method, path, timestamp, nonce, bodyBytes)

	var lastErr error
	for attempt := 0; attempt <= MaxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retry with exponential backoff
			wait := time.Duration(attempt) * RetryWaitMin
			if wait > RetryWaitMax {
				wait = RetryWaitMax
			}
			time.Sleep(wait)

			// Regenerate timestamp and nonce for retry
			timestamp = time.Now().Unix()
			nonce, _ = uuid.GenerateUUID()
			signature = c.generateHMAC(method, path, timestamp, nonce, bodyBytes)
		}

		// Create request
		var bodyReader io.Reader
		if bodyBytes != nil {
			bodyReader = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Set headers
		req.Header.Set("Authorization", c.buildAuthorizationHeader(signature, nonce, timestamp))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", c.UserAgent)

		// Perform request
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			continue
		}

		// Read response body
		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			continue
		}

		// Check for rate limiting (429)
		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := resp.Header.Get("Retry-After")
			if retryAfter != "" {
				if seconds, err := strconv.Atoi(retryAfter); err == nil {
					time.Sleep(time.Duration(seconds) * time.Second)
				}
			}
			lastErr = &APIError{
				StatusCode: resp.StatusCode,
				Message:    "rate limit exceeded",
			}
			continue
		}

		// Check for server errors (5xx) - retry
		if resp.StatusCode >= 500 {
			lastErr = &APIError{
				StatusCode: resp.StatusCode,
				Message:    string(respBody),
			}
			continue
		}

		// Check for client errors (4xx) - don't retry
		if resp.StatusCode >= 400 {
			apiErr := &APIError{
				StatusCode: resp.StatusCode,
				Message:    string(respBody),
			}

			// Try to parse validation errors
			var badRequest struct {
				ValidationErrors []ValidationError `json:"validation_errors"`
			}
			if json.Unmarshal(respBody, &badRequest) == nil && len(badRequest.ValidationErrors) > 0 {
				apiErr.ValidationErrors = badRequest.ValidationErrors
			}

			return resp, respBody, apiErr
		}

		// Success
		return resp, respBody, nil
	}

	return nil, nil, lastErr
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, path string) (*http.Response, []byte, error) {
	return c.doRequest(ctx, http.MethodGet, path, nil)
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, path string, body interface{}) (*http.Response, []byte, error) {
	return c.doRequest(ctx, http.MethodPost, path, body)
}

// Put performs a PUT request
func (c *Client) Put(ctx context.Context, path string, body interface{}) (*http.Response, []byte, error) {
	return c.doRequest(ctx, http.MethodPut, path, body)
}

// Delete performs a DELETE request
func (c *Client) Delete(ctx context.Context, path string) (*http.Response, []byte, error) {
	return c.doRequest(ctx, http.MethodDelete, path, nil)
}

// GetRateLimitInfo extracts rate limit information from response headers
func GetRateLimitInfo(resp *http.Response) *RateLimitInfo {
	if resp == nil {
		return nil
	}

	info := &RateLimitInfo{}

	if limit := resp.Header.Get("X-RateLimit-Limit"); limit != "" {
		info.Limit, _ = strconv.Atoi(limit)
	}
	if usage := resp.Header.Get("X-RateLimit-Usage"); usage != "" {
		info.Usage, _ = strconv.Atoi(usage)
	}
	if remaining := resp.Header.Get("X-RateLimit-Remaining"); remaining != "" {
		info.Remaining, _ = strconv.Atoi(remaining)
	}
	if reset := resp.Header.Get("X-RateLimit-Reset"); reset != "" {
		info.Reset, _ = strconv.Atoi(reset)
	}

	return info
}

// GetLocationHeader extracts the Location header from a response
func GetLocationHeader(resp *http.Response) string {
	if resp == nil {
		return ""
	}
	return resp.Header.Get("Location")
}

// GetPagingInfo extracts paging information from response headers
type PagingInfo struct {
	Skipped      int
	Take         int
	TotalResults int
}

// GetPagingInfo extracts paging information from response headers
func GetPagingInfo(resp *http.Response) *PagingInfo {
	if resp == nil {
		return nil
	}

	info := &PagingInfo{}

	if skipped := resp.Header.Get("X-Paging-Skipped"); skipped != "" {
		info.Skipped, _ = strconv.Atoi(skipped)
	}
	if take := resp.Header.Get("X-Paging-Take"); take != "" {
		info.Take, _ = strconv.Atoi(take)
	}
	if totalResults := resp.Header.Get("X-Paging-TotalResults"); totalResults != "" {
		info.TotalResults, _ = strconv.Atoi(totalResults)
	}

	return info
}
