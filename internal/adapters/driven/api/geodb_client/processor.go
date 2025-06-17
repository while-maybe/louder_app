package geodbclient

import (
	"context"
	"log"
	"net/url"
)

type processorResult struct {
	response *GeoDBAPIResponse
	err      error
}

type processor struct {
	httpClient *httpClient
}

func newProcessor(client *httpClient) *processor {
	return &processor{
		httpClient: client,
	}
}

// execute performs the API call asynchronously and send the result to a channel
func (p *processor) execute(ctx context.Context, endpoint string, params url.Values) <-chan processorResult {
	resultChan := make(chan processorResult, 1)

	go func() {
		defer close(resultChan)
		log.Printf("Processor: Goroutine starting API call to endpoint '%s' with params: %v", endpoint, params)

		response, err := p.httpClient.queryAPI(ctx, endpoint, params)
		// I won't handle the error here, instead send that in the channel for someone else to deal with

		resultChan <- processorResult{response: response, err: err}
	}()

	return resultChan
}
