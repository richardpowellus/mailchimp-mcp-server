package mailchimp

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// APIError represents a Mailchimp application/problem+json error.
type APIError struct {
	StatusCode int    `json:"status"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	Detail     string `json:"detail"`
	Instance   string `json:"instance"`
}

func (e *APIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("Mailchimp %d %s: %s", e.StatusCode, e.Title, e.Detail)
	}
	return fmt.Sprintf("Mailchimp %d: %s", e.StatusCode, e.Title)
}

func (e *APIError) Is429() bool { return e.StatusCode == 429 }
func (e *APIError) Is404() bool { return e.StatusCode == 404 }

// Client wraps the Mailchimp Marketing API v3.
type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client
	sem     chan struct{} // concurrency limiter (10 simultaneous connections)
}

// NewClient creates a new Mailchimp API client.
// The API key format is "key-datacenter" (e.g., "abc123-us2").
// Panics if the key format is invalid.
func NewClient(apiKey string) *Client {
	idx := strings.LastIndex(apiKey, "-")
	if idx < 0 || idx == len(apiKey)-1 {
		panic(fmt.Sprintf("invalid Mailchimp API key format: expected 'key-datacenter' (e.g., 'abc123-us2'), got key ending in %q", apiKey[max(0, len(apiKey)-5):]))
	}
	dc := apiKey[idx+1:]
	return &Client{
		apiKey:  apiKey,
		baseURL: fmt.Sprintf("https://%s.api.mailchimp.com/3.0", dc),
		http: &http.Client{
			Timeout: 120 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        20,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		sem: make(chan struct{}, 10), // Mailchimp limit: 10 concurrent connections
	}
}

// authHeader returns the Basic Auth header value.
func (c *Client) authHeader() string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte("anystring:"+c.apiKey))
}

// SubscriberHash returns the MD5 hash of a lowercase email address,
// as required by all Mailchimp member endpoints.
func SubscriberHash(email string) string {
	h := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(email))))
	return fmt.Sprintf("%x", h)
}

// doRequest executes an HTTP request with concurrency limiting, retry, and error parsing.
func (c *Client) doRequest(ctx context.Context, method, path string, body []byte) (json.RawMessage, error) {
	// Acquire semaphore
	select {
	case c.sem <- struct{}{}:
		defer func() { <-c.sem }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	isRead := method == "GET"
	maxRetries := 3
	backoff := []time.Duration{1 * time.Second, 2 * time.Second, 4 * time.Second}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			wait := backoff[min(attempt-1, len(backoff)-1)]
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(wait):
			}
		}

		var bodyReader io.Reader
		if body != nil {
			bodyReader = strings.NewReader(string(body))
		}

		req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", c.authHeader())
		req.Header.Set("Accept", "application/json")
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := c.http.Do(req)
		if err != nil {
			lastErr = err
			if isRead {
				continue
			}
			return nil, err
		}

		respBody, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("reading response: %w", err)
		}

		// Success: 2xx
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Handle 204 No Content
			if resp.StatusCode == 204 || len(respBody) == 0 {
				return json.RawMessage(`{"success":true}`), nil
			}
			// Strip _links from response
			cleaned := stripLinks(respBody)
			return cleaned, nil
		}

		// Parse error
		apiErr := &APIError{StatusCode: resp.StatusCode}
		_ = json.Unmarshal(respBody, apiErr)
		if apiErr.Title == "" {
			apiErr.Title = http.StatusText(resp.StatusCode)
		}

		// Retry logic: GET retries on transient errors; writes only on 429
		if resp.StatusCode == 429 {
			lastErr = apiErr
			// Honor Retry-After header if present
			if ra := resp.Header.Get("Retry-After"); ra != "" {
				if secs, err := strconv.Atoi(ra); err == nil && secs > 0 && secs <= 300 {
					select {
					case <-ctx.Done():
						return nil, ctx.Err()
					case <-time.After(time.Duration(secs) * time.Second):
					}
					continue
				}
			}
			continue
		}
		if isRead && (resp.StatusCode >= 500) {
			lastErr = apiErr
			continue
		}

		return nil, apiErr
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", maxRetries, lastErr)
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, path string) (json.RawMessage, error) {
	return c.doRequest(ctx, "GET", path, nil)
}

// Post performs a POST request with a JSON body.
func (c *Client) Post(ctx context.Context, path string, body any) (json.RawMessage, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return c.doRequest(ctx, "POST", path, data)
}

// Put performs a PUT request with a JSON body.
func (c *Client) Put(ctx context.Context, path string, body any) (json.RawMessage, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return c.doRequest(ctx, "PUT", path, data)
}

// Patch performs a PATCH request with a JSON body.
func (c *Client) Patch(ctx context.Context, path string, body any) (json.RawMessage, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return c.doRequest(ctx, "PATCH", path, data)
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, path string) (json.RawMessage, error) {
	return c.doRequest(ctx, "DELETE", path, nil)
}

// PostRaw performs a POST request with a raw JSON body (for json.RawMessage passthrough).
func (c *Client) PostRaw(ctx context.Context, path string, body json.RawMessage) (json.RawMessage, error) {
	return c.doRequest(ctx, "POST", path, body)
}

// PutRaw performs a PUT request with a raw JSON body.
func (c *Client) PutRaw(ctx context.Context, path string, body json.RawMessage) (json.RawMessage, error) {
	return c.doRequest(ctx, "PUT", path, body)
}

// PatchRaw performs a PATCH request with a raw JSON body.
func (c *Client) PatchRaw(ctx context.Context, path string, body json.RawMessage) (json.RawMessage, error) {
	return c.doRequest(ctx, "PATCH", path, body)
}

// FetchAll auto-paginates a Mailchimp list endpoint using offset/count.
// It collects all items from the specified key in the response.
// For example, FetchAll(ctx, "/campaigns", "campaigns") fetches all campaigns.
func (c *Client) FetchAll(ctx context.Context, path, itemsKey string) ([]json.RawMessage, error) {
	var allItems []json.RawMessage
	offset := 0
	count := 100 // fetch 100 per page for efficiency

	sep := "?"
	if strings.Contains(path, "?") {
		sep = "&"
	}

	for {
		pagePath := fmt.Sprintf("%s%soffset=%d&count=%d", path, sep, offset, count)
		raw, err := c.Get(ctx, pagePath)
		if err != nil {
			return nil, err
		}

		var wrapper map[string]json.RawMessage
		if err := json.Unmarshal(raw, &wrapper); err != nil {
			return nil, fmt.Errorf("parsing paginated response: %w", err)
		}

		itemsRaw, ok := wrapper[itemsKey]
		if !ok {
			break
		}

		var items []json.RawMessage
		if err := json.Unmarshal(itemsRaw, &items); err != nil {
			return nil, fmt.Errorf("parsing items array: %w", err)
		}

		allItems = append(allItems, items...)

		// Check total_items to know when we're done
		var meta struct {
			TotalItems int `json:"total_items"`
		}
		_ = json.Unmarshal(raw, &meta)

		if meta.TotalItems > 0 && len(allItems) >= meta.TotalItems {
			break
		}
		if len(items) < count {
			break
		}

		offset += count
	}

	return allItems, nil
}

// stripLinks removes _links arrays from JSON responses recursively.
// Mailchimp includes HATEOAS links in every response which waste tokens.
func stripLinks(data []byte) json.RawMessage {
	var obj any
	if err := json.Unmarshal(data, &obj); err != nil {
		return data // return as-is if not valid JSON
	}
	cleaned := stripLinksRecursive(obj)
	result, err := json.Marshal(cleaned)
	if err != nil {
		return data
	}
	return result
}

func stripLinksRecursive(v any) any {
	switch val := v.(type) {
	case map[string]any:
		result := make(map[string]any, len(val))
		for k, v := range val {
			if k == "_links" {
				continue
			}
			result[k] = stripLinksRecursive(v)
		}
		return result
	case []any:
		result := make([]any, len(val))
		for i, v := range val {
			result[i] = stripLinksRecursive(v)
		}
		return result
	default:
		return v
	}
}
