package overpass

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"

	"github.com/Otherotter/city-explorer/services/collector/internal/models"
)

const overpassURL = "https://overpass-api.de/api/interpreter"

type Element struct {
	Type string  `json:"type"`
	ID   int64   `json:"id"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Tags struct {
		Name        string `json:"name"`
		Amenity     string `json:"amenity"`
		Cuisine     string `json:"cuisine"`
		Website     string `json:"website"`
		Phone       string `json:"phone"`
		Street      string `json:"addr:street"`
		HouseNumber string `json:"addr:housenumber"`
		City        string `json:"addr:city"`
	} `json:"tags"`
}

type Response struct {
	Elements []Element `json:"elements"`
}

// FetchPlaces calls the Overpass API and returns
// normalized Place models ready for the database.
// ctx is passed through so traces span across
// the full request journey.
func FetchPlaces(ctx context.Context, cityID int, categoryID int, amenities []string, bbox string) ([]models.Place, error) {
	query := buildQuery(amenities, bbox)

	resp, err := callOverpass(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("overpass api call failed: %w", err)
	}

	places := normalize(resp.Elements, cityID, categoryID)
	return places, nil
}

func buildQuery(amenities []string, bbox string) string {
	var sb strings.Builder
	sb.WriteString("[out:json][timeout:25];(\n")
	for _, amenity := range amenities {
		sb.WriteString(fmt.Sprintf(
			"  node[\"amenity\"=\"%s\"](%s);\n",
			amenity, bbox,
		))
	}
	sb.WriteString(");out body;")
	return sb.String()
}

// callOverpass makes the HTTP request to the Overpass API.
// It creates a trace span so failures are visible in Tempo.
func callOverpass(ctx context.Context, query string) (*Response, error) {
	start := time.Now()

	// Start a trace span for this external API call
	// This is the key instrumentation — when this fails
	// the span shows as red in Tempo immediately
	tracer := otel.Tracer("city-explorer-collector")
	ctx, span := tracer.Start(ctx, "overpass.http_request")
	defer span.End()

	// Attach context to the span so it is searchable
	span.SetAttributes(
		attribute.String("url", overpassURL),
		attribute.Int("query.length", len(query)),
	)

	slog.Info("calling overpass api",
		"url", overpassURL,
		"query_length", len(query),
	)

	formData := url.Values{"data": {query}}
	encodedBody := formData.Encode()

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		overpassURL,
		strings.NewReader(encodedBody),
	)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to build request")
		span.SetAttributes(attribute.String("error", err.Error()))
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "city-explorer/1.0 (personal project)")

	resp, err := client.Do(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "http request failed")
		span.SetAttributes(attribute.String("error", err.Error()))
		slog.Error("overpass http request failed",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	duration := time.Since(start)

	// Always record the status code on the span
	// This is what makes failures visible in Tempo
	span.SetAttributes(
		attribute.Int("http.status_code", resp.StatusCode),
		attribute.Int64("duration_ms", duration.Milliseconds()),
	)

	if resp.StatusCode != http.StatusOK {
		body := make([]byte, 500)
		n, _ := resp.Body.Read(body)

		// Mark span as error — this turns it RED in Tempo
		span.RecordError(err)
		span.SetStatus(codes.Error, fmt.Sprintf("overpass returned %d", resp.StatusCode))

		slog.Error("overpass returned error status",
			"status_code", resp.StatusCode,
			"response_body", string(body[:n]),
			"duration_ms", duration.Milliseconds(),
		)
		return nil, fmt.Errorf("overpass returned status %d", resp.StatusCode)
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to decode response")
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Success — mark span with result count
	span.SetAttributes(attribute.Int("elements.returned", len(result.Elements)))
	span.SetStatus(codes.Ok, "")

	slog.Info("overpass request complete",
		"status_code", resp.StatusCode,
		"elements", len(result.Elements),
		"duration_ms", duration.Milliseconds(),
	)

	return &result, nil
}

func normalize(elements []Element, cityID int, categoryID int) []models.Place {
	places := make([]models.Place, 0, len(elements))

	for _, el := range elements {
		if el.Tags.Name == "" {
			continue
		}

		place := models.Place{
			CityID:      cityID,
			CategoryID:  categoryID,
			Name:        el.Tags.Name,
			Subcategory: el.Tags.Amenity,
			Latitude:    el.Lat,
			Longitude:   el.Lon,
			Website:     el.Tags.Website,
			Phone:       el.Tags.Phone,
			Source:      "overpass",
			SourceID:    strconv.FormatInt(el.ID, 10),
		}

		if el.Tags.HouseNumber != "" && el.Tags.Street != "" {
			place.Address = fmt.Sprintf("%s %s",
				el.Tags.HouseNumber,
				el.Tags.Street,
			)
		}

		if el.Tags.Cuisine != "" {
			place.Subcategory = el.Tags.Cuisine
		}

		places = append(places, place)
	}

	return places
}
