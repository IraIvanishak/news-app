package unsplash

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/viper"
)

// Client is the HTTP client for Unsplash API
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	appName    string
}

// NewUnsplashClient creates a new Unsplash API client
func NewUnsplashClient() *Client {
	return &Client{
		apiKey:  viper.GetString("UNSPLASH_ACCESS_KEY"),
		baseURL: "https://api.unsplash.com",
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		appName: viper.GetString("UNSPLASH_APP_NAME"),
	}
}

// SearchPhotos searches for photos on Unsplash
func (c *Client) SearchPhotos(ctx context.Context, query string, orientation string, perPage int) (*SearchResponse, error) {
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	// Build URL with query parameters
	u, err := url.Parse(fmt.Sprintf("%s/search/photos", c.baseURL))
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	params := url.Values{}
	params.Add("query", query)
	params.Add("per_page", fmt.Sprintf("%d", perPage))
	if orientation != "" {
		params.Add("orientation", orientation)
	}
	u.RawQuery = params.Encode()

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Client-ID %s", c.apiKey))
	req.Header.Set("Accept", "application/json")

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Handle error responses
	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return nil, fmt.Errorf("invalid Unsplash API key")
		case http.StatusTooManyRequests:
			return nil, fmt.Errorf("rate limit exceeded, please try again later")
		case http.StatusNotFound:
			return nil, fmt.Errorf("endpoint not found")
		default:
			return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		}
	}

	// Parse response
	var searchResponse SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &searchResponse, nil
}

// TriggerDownload triggers the download endpoint to track photo usage
func (c *Client) TriggerDownload(ctx context.Context, downloadLocation string) error {
	if downloadLocation == "" {
		return fmt.Errorf("download location cannot be empty")
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", downloadLocation, nil)
	if err != nil {
		return fmt.Errorf("failed to create download request: %w", err)
	}

	// Set authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Client-ID %s", c.apiKey))

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to trigger download: %w", err)
	}
	defer resp.Body.Close()

	// We don't fail the operation if download tracking fails
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download tracking returned status: %d", resp.StatusCode)
	}

	return nil
}
