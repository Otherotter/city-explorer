package collector

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// Client calls the Collector service.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a Collector client.
// The base URL comes from the environment.
func NewClient() *Client {
	baseURL := os.Getenv("COLLECTOR_URL")
	if baseURL == "" {
		baseURL = "http://collector:8081"
	}

	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Collect triggers data collection for a city and category.
// It blocks until the Collector finishes.
// The trace context is passed through so Tempo shows
// the collection as part of the original request trace.
func (c *Client) Collect(ctx context.Context, city string, category string) error {
	tracer := otel.Tracer("city-explorer-api")

	ctx, span := tracer.Start(ctx, "collector.trigger")
	span.SetAttributes(
		attribute.String("city", city),
		attribute.String("category", category),
	)
	defer span.End()

	url := fmt.Sprintf("%s/collect/%s/%s", c.baseURL, city, category)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("collector request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("collector returned status %d", resp.StatusCode)
	}

	return nil
}
