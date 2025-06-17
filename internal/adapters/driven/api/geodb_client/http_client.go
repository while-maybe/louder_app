package geodbclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type httpClient struct {
	baseURL          string
	apiKeyHeaderName string
	apiKeyValue      string
	httpClient       *http.Client
}

// baseURL: "https://wft-geo-db.p.rapidapi.com",

func NewHTTPClent(baseURL, apiKey string) *httpClient {
	return &httpClient{
		baseURL:          baseURL,
		apiKeyHeaderName: "x-rapidapi-key",
		apiKeyValue:      apiKey,

		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// endpoint := "/v1/geo/countries"

func (c *httpClient) queryAPI(ctx context.Context, endpoint string, params url.Values) (*GeoDBAPIResponse, error) {

	// join base url and endpoint and check for errors
	joinedURL, err := url.JoinPath(c.baseURL, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to join base URL and endpoint: %w", err)
	}

	// parse url and check for errors
	parsedURL, err := url.Parse(joinedURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	if params != nil {
		// encode the params into the query
		parsedURL.RawQuery = params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	log.Printf("HTTP Client: Sending request to %s", parsedURL.String())

	// add the auth key to the header
	req.Header.Set(c.apiKeyHeaderName, c.apiKeyValue)

	// use the http client injected to the function
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// context cancelled or deadline?
		switch {
		case errors.Is(err, context.Canceled):
			log.Printf("HTTP Client: Request cancelled: %v", err)
			return nil, context.Canceled
		case errors.Is(err, context.DeadlineExceeded):
			log.Printf("HTTP Client: Request timed out: %v", err)
			return nil, context.DeadlineExceeded
		default:
			return nil, fmt.Errorf("http_client: httpClient.Do: %w", err)
		}
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with code %s: %s", resp.Status, string(bodyBytes))
	}

	// read and parse the body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp GeoDBAPIResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		log.Printf("HTTP Client: Failed to unmarshal JSON. Body was: %s", string(bodyBytes))
		return nil, fmt.Errorf("failed to unmarshall json data: %w", err)
	}

	return &apiResp, nil
}
