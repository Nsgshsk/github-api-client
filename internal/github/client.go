package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseURL    = "https://api.github.com"
	userAgent  = "app/github-api-client"
	apiVersion = "2022-11-28"

	defaultTimeout   = 30 * time.Second
	maxResponseBytes = 10 << 20
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

type Option func(*Client)

func WithToken(token string) Option {
	return func(c *Client) {
		c.token = strings.TrimSpace(token)
	}
}

func WithBaseURL(rawURL string) Option {
	return func(c *Client) {
		c.baseURL = strings.TrimSuffix(rawURL, "/")
	}
}

func NewClient(opts ...Option) *Client {
	client := &Client{
		httpClient: &http.Client{Timeout: defaultTimeout},
		baseURL:    baseURL,
	}
	for _, opt := range opts {
		opt(client)
	}
	return client
}

// APIError describes a non-successful response returned by the GitHub API.
type APIError struct {
	StatusCode int
	Message    string
	URL        string
}

func (e *APIError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("github: %s returned status %d", e.URL, e.StatusCode)
	}
	return fmt.Sprintf("github: %s returned status %d: %s", e.URL, e.StatusCode, e.Message)
}

// IsRateLimited reports whether err was caused by GitHub's rate limiting.
func IsRateLimited(err error) bool {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return false
	}
	if apiErr.StatusCode == http.StatusTooManyRequests {
		return true
	}
	return apiErr.StatusCode == http.StatusForbidden &&
		strings.Contains(strings.ToLower(apiErr.Message), "rate limit")
}

func (c *Client) do(path string, out any) (next string, err error) {
	requestURL := path
	if !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") {
		requestURL = c.baseURL + path
	}

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", apiVersion)
	req.Header.Set("User-Agent", userAgent)
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request %s: %w", requestURL, err)
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close response body: %w", closeErr)
		}
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, readErr := io.ReadAll(io.LimitReader(resp.Body, 4096))
		if readErr != nil {
			return "", fmt.Errorf("read error response from %s: %w", requestURL, readErr)
		}
		return "", newAPIError(resp.StatusCode, body, requestURL)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes+1))
	if err != nil {
		return "", fmt.Errorf("read response from %s: %w", requestURL, err)
	}
	if len(body) > maxResponseBytes {
		return "", fmt.Errorf("response from %s exceeds %d bytes", requestURL, maxResponseBytes)
	}

	if out != nil {
		if err = json.Unmarshal(body, out); err != nil {
			return "", fmt.Errorf("decode response from %s: %w", requestURL, err)
		}
	}

	return nextPageURL(resp.Header.Get("Link")), nil
}

func newAPIError(statusCode int, body []byte, requestURL string) error {
	var payload struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &payload); err == nil && payload.Message != "" {
		return &APIError{StatusCode: statusCode, Message: payload.Message, URL: requestURL}
	}
	return &APIError{StatusCode: statusCode, Message: strings.TrimSpace(string(body)), URL: requestURL}
}

func getPaged[T any](c *Client, path string) ([]T, error) {
	var all []T
	next := path
	seen := make(map[string]struct{})
	for next != "" {
		if _, visited := seen[next]; visited {
			return nil, fmt.Errorf("pagination loop detected at %s", next)
		}
		seen[next] = struct{}{}

		var page []T
		nextURL, err := c.do(next, &page)
		if err != nil {
			return nil, err
		}
		all = append(all, page...)
		next = nextURL
	}
	return all, nil
}

func nextPageURL(link string) string {
	for _, part := range strings.Split(link, ",") {
		part = strings.TrimSpace(part)
		if !strings.Contains(part, `rel="next"`) {
			continue
		}
		start := strings.Index(part, "<")
		end := strings.Index(part, ">")
		if start >= 0 && end > start {
			return part[start+1 : end]
		}
	}
	return ""
}

func escapeRepo(fullName string) string {
	owner, repo, found := strings.Cut(fullName, "/")
	if !found {
		return url.PathEscape(fullName)
	}
	return url.PathEscape(owner) + "/" + url.PathEscape(repo)
}
