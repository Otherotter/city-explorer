package overpass

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Otherotter/city-explorer/services/collector/internal/models"
)

// The Overpass API endpoint. Public and free.
const overpassURL = "https://overpass-api.de/api/interpreter"

// Element is a single result from the Overpass API.
// OSM returns nodes (points), ways (lines/areas),
// and relations. We only care about nodes for now.
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

// Response is the full Overpass API response.
type Response struct {
	Elements []Element `json:"elements"`
}

// FetchPlaces calls the Overpass API and returns
// normalized Place models ready for the database.
//
// amenities is a list of OSM amenity types.
// Example: []string{"restaurant", "cafe", "library"}
//
// bbox is the bounding box for the city.
// Format: south, west, north, east
func FetchPlaces(cityID int, categoryID int, amenities []string, bbox string) ([]models.Place, error) {
	// Build the Overpass QL query
	// This asks for all nodes within the
	// bounding box that match our amenity types
	query := buildQuery(amenities, bbox)

	// Call the API
	resp, err := callOverpass(query)
	if err != nil {
		return nil, fmt.Errorf("overpass api call failed: %w", err)
	}

	// Normalize each element into our Place model
	places := normalize(resp.Elements, cityID, categoryID)

	return places, nil
}

// buildQuery constructs an Overpass QL query string.
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
func callOverpass(query string) (*Response, error) {

	// log.Printf("[callOverpass] sending query:\n%s", query)
	// log.Printf("[callOverpass] query length: %d", len(query))
	// log.Printf("[callOverpass] url: %s", overpassURL)
	slog.Info("calling OverPass APIs")
	slog.Debug("query strucuture",
		"query", query,
		"query len", len(query),
		"url", overpassURL,
	)

	// Build the form body manually so we can inspect it
	formData := url.Values{"data": {query}}
	encodedBody := formData.Encode()

	log.Printf("[callOverpass] encoded body length: %d bytes", len(encodedBody))
	log.Printf("[callOverpass] encoded body preview: %.100s...", encodedBody)
	slog.Debug("encoded body",
		"encoded body length", len(encodedBody),
		"encoded body preview", encodedBody,
	)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Build request manually instead of PostForm
	// so we can set and inspect every header
	req, err := http.NewRequest(
		http.MethodPost,
		overpassURL,
		strings.NewReader(encodedBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	// This is what PostForm sets internally
	// We set it explicitly so we can confirm it is there
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Adding a User-Agent is good practice for public APIs
	// Some APIs block requests with no User-Agent
	req.Header.Set("User-Agent", "city-explorer/1.0 (personal project)")

	// log.Printf("[callOverpass] request headers: %v", req.Header)
	// log.Printf("[callOverpass] sending to: %s", overpassURL)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	// log.Printf("[callOverpass] response status: %d", resp.StatusCode)
	// log.Printf("[callOverpass] response headers: %v", resp.Header)
	slog.Debug("built format",
		"request_headers", req.Header,
		"overpassURL", overpassURL,
		"response_status", resp.StatusCode,
		"response_headers", resp.Header,
	)

	if resp.StatusCode != http.StatusOK {
		// Read full error body from Overpass
		body := make([]byte, 1000)
		n, _ := resp.Body.Read(body)
		log.Printf("[callOverpass] error response body:\n%s", string(body[:n]))
		return nil, fmt.Errorf("overpass returned status %d", resp.StatusCode)
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// log.Printf("[callOverpass] success: %d elements returned", len(result.Elements))

	return &result, nil
	///----
	//curl -X POST "https://overpass-api.de/api/interpreter" --data 'data=[out:json][timeout:25];(node["amenity"="restaurant"](40.4774,-74.2591,40.9176,-73.7004););out body;'  -v

	// resp, err := client.PostForm(overpassURL, url.Values{
	// 	"data": {query},
	// })

	// if err != nil {
	// 	return nil, fmt.Errorf("http request failed: %w", err)
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	return nil, fmt.Errorf("overpass returned status %d", resp.StatusCode)
	// }

	// var result Response
	// if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
	// 	return nil, fmt.Errorf("failed to decode response: %w", err)
	// }

	// return &result, nil
}

// normalize converts raw Overpass elements into
// your internal Place model. This is where you
// own the data shape — not the external API.
func normalize(elements []Element, cityID int, categoryID int) []models.Place {
	places := make([]models.Place, 0, len(elements))

	for _, el := range elements {
		// Skip elements with no name
		// They are not useful to us
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

		// Build address from parts if available
		if el.Tags.HouseNumber != "" && el.Tags.Street != "" {
			place.Address = fmt.Sprintf(
				"%s %s",
				el.Tags.HouseNumber,
				el.Tags.Street,
			)
		}

		// Use cuisine as subcategory for food places
		if el.Tags.Cuisine != "" {
			place.Subcategory = el.Tags.Cuisine
		}

		places = append(places, place)
	}

	return places
}
