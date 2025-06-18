package geodbclient

import (
	"context"
	"errors"
	"fmt"
	"log"
	"louder/internal/core/domain"
	"louder/internal/core/service/countrycore"
	"louder/internal/core/service/currencycore"
	"time"
)

type Provider struct {
	httpClient         *httpClient
	currencyRepo       currencycore.Repository
	apiCountryEndpoint string
	apiPageLimit       int
	apiRateLimitSleep  time.Duration
}

func NewProvider(baseURL, countryEndpoint, apiKey string, currencyRepo currencycore.Repository, pageLimit int, rateLimitSleep time.Duration) countrycore.ExternalCountryProvider {
	client := NewHTTPClient(baseURL, apiKey)
	return &Provider{
		httpClient:         client,
		currencyRepo:       currencyRepo,
		apiCountryEndpoint: countryEndpoint,
		apiPageLimit:       pageLimit,
		apiRateLimitSleep:  rateLimitSleep,
	}
}

func (p *Provider) FetchAllCountries(ctx context.Context) ([]*domain.Country, error) {
	log.Println("GeoDB Provider: FetchAllCountries called.")
	var domainCountries []*domain.Country

	// create a processor
	proc := newProcessor(p.httpClient)

	// create a paginator
	countriesPaginator := NewPaginator(proc, p.apiCountryEndpoint, p.apiPageLimit)

	for countriesPaginator.hasNext {
		select {
		case <-ctx.Done():
			return domainCountries, fmt.Errorf("FetchAllCountries aborted by context: %w", ctx.Err())
		default:
		}

		countryDTOs, err := countriesPaginator.NextPage(ctx)
		if err != nil {
			return domainCountries, fmt.Errorf("error fetching page via paginator: %w", err)
		}

		// if countryDTOs is empty and there are no more pages
		if len(countryDTOs) == 0 && !countriesPaginator.HasNext() {
			break
		}

		for _, dto := range countryDTOs {
			select {
			case <-ctx.Done():
				return domainCountries, fmt.Errorf("FetchAllCountries aborted during DTO mapping: %w", ctx.Err())

			default:
			}

			domainCountry, mapErr := p.mapDTOToDomainCountry(ctx, dto)
			if mapErr != nil {
				log.Printf("ERROR Provider: Failed to map country DTO %s: %v. Skipping.", dto.CountryCode, mapErr)
				continue
			}

			domainCountries = append(domainCountries, domainCountry)
		}

		if countriesPaginator.HasNext() {
			log.Printf("ProviderImpl: Fetched page, %d countries total so far. Sleeping for rate limit.", len(domainCountries))

			select {
			case <-time.After(p.apiRateLimitSleep):
			case <-ctx.Done():
				return domainCountries, fmt.Errorf("FetchAllCountries aborted during rate limit sleep: %w", ctx.Err())
			}
		}
	}
	log.Printf("ProviderImpl: Successfully fetched and mapped %d countries.", len(domainCountries))
	return domainCountries, nil
}

func (p *Provider) GetTotalCountryCountFromAPI(ctx context.Context) (int, error) {
	proc := newProcessor(p.httpClient)
	pg := NewPaginator(proc, p.apiCountryEndpoint, p.apiPageLimit)

	_, err := pg.NextPage(ctx)
	if err != nil && pg.TotalCount() == -1 {
		return 0, fmt.Errorf("GetTotalCountryCountFromAPI: initial page fetch failed: %w", err)
	}

	if pg.TotalCount() == -1 {
		return 0, errors.New("GetTotalCountryCountFromAPI: total count not determined")
	}

	return pg.TotalCount(), nil
}
