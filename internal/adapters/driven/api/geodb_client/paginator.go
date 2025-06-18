package geodbclient

import (
	"context"
	"errors"
	"fmt"
	"log"
	"louder/internal/adapters/driven/db/dbcommon"
	"louder/internal/core/domain"
	"net/url"
	"strconv"
)

type paginator struct {
	proc       *processor
	endpoint   string
	offset     int
	limit      int
	totalCount int
	hasNext    bool
}

func NewPaginator(proc *processor, endpoint string, limitPerPage int) *paginator {
	return &paginator{
		proc:       proc,
		endpoint:   endpoint,
		offset:     0,
		limit:      limitPerPage,
		totalCount: -1,
		hasNext:    true,
	}
}

func (p *paginator) HasNext() bool {
	return p.hasNext
}

func (p *paginator) TotalCount() int {
	return p.totalCount
}

func (p *paginator) NextPage(ctx context.Context) ([]CountryDTO, error) {
	if !p.HasNext() {
		return nil, nil // no pages, no errors
	}

	params := url.Values{}
	params.Set("limit", strconv.Itoa(p.limit))
	params.Add("offset", strconv.Itoa(p.offset))
	log.Printf("Paginator: Requesting next page - endpoint: %s, offset: %d, limit: %d", p.endpoint, p.offset, p.limit)

	resultChan := p.proc.execute(ctx, p.endpoint, params)

	var apiResponse *GeoDBAPIResponse
	var pageErr error

	select {
	case result := <-resultChan:
		switch {
		case result.err != nil:
			pageErr = fmt.Errorf("paginator: API call failed for offset %d (endpoint %s): %w", p.offset, p.endpoint, result.err)
			p.hasNext = false

		default:
			apiResponse = result.response

			// if p.totalCount is -1, we didn't know this value before
			if p.totalCount == -1 && apiResponse != nil {
				p.totalCount = apiResponse.Metadata.Count
			}
			log.Printf("Paginator: Discovered total API count for %s: %d", p.endpoint, p.totalCount)

			switch {
			// an empty response or no countries received means there are no next pages
			case apiResponse == nil || len(apiResponse.Countries) == 0:
				p.hasNext = false

			default:
				// increase the offset by the right amount
				p.offset += len(apiResponse.Countries)
				// if this the first time we get and we've got the total count already (1 page only?)
				if p.totalCount != -1 && p.offset >= p.totalCount {
					p.hasNext = false
				}
			}
		}

	// ctx cancelled?
	case <-ctx.Done():
		pageErr = fmt.Errorf("paginator: context cancelled for offset %d (endpoint %s): %w", p.offset, p.endpoint, ctx.Err())
		p.hasNext = false
	}

	// if we could not get the page
	if pageErr != nil {
		return nil, pageErr
	}

	// page was empty?
	if apiResponse == nil {
		return []CountryDTO{}, nil
	}

	return apiResponse.Countries, nil
}

// mapDTOToDomainCountry converts an API DTO (CountryDTO) to a domain.Country object
func (p *Provider) mapDTOToDomainCountry(ctx context.Context, dto CountryDTO) (*domain.Country, error) {
	if dto.CountryCode == "" {
		return nil, errors.New("mapDTO: API DTO has empty country code")
	}

	if dto.CountryName == "" {
		return nil, errors.New("mapDTO: API DTO has empty country name")
	}

	countryCode, err := domain.NewCountryCode(dto.CountryCode)
	if err != nil {
		return nil, fmt.Errorf("mapDTO: invalid country code '%s' from API: %w", dto.CountryCode, err)
	}

	// wikidataid can be null. For now...
	var wikiID domain.WikiCode
	if dto.WikiDataId != "" {
		wikiID, err = domain.NewWikiCode(dto.WikiDataId)

		if err != nil {
			log.Printf("WARN mapDTO: Invalid WikiDataID format '%s' for country %s. Using empty. Error: %v", dto.WikiDataId, dto.CountryCode, err)
			wikiID = domain.WikiCode("")
		}
	}

	var domainCurrencies []domain.Currency
	for _, currencyCodeStr := range dto.CurrencyCodes {
		if currencyCodeStr == "" {
			log.Printf("WARN mapDTO: Empty currency code string received for country %s. Skipping.", dto.CountryCode)
			continue
		}

		cc, err := domain.NewCurrencyCode(currencyCodeStr)
		if err != nil {
			log.Printf("WARN mapDTO: Invalid currency code string '%s' for country %s. Skipping this currency. Error: %v", currencyCodeStr, dto.CountryCode, err)
			continue
		}

		// get the currency from the DB if exists
		currency, err := p.currencyRepo.GetByID(ctx, cc)
		if errors.Is(err, dbcommon.ErrNotFound) {
			log.Printf("WARN mapDTO: Currency %s for country %s not in local DB. Creating placeholder domain object for now.", cc.String(), dto.CountryCode)

			placeholderName := fmt.Sprintf("Currency %s (Auto-from API sync)", cc.String())

			newCurrency, ncErr := domain.NewCurrency(cc, placeholderName)
			if ncErr != nil {
				log.Printf("ERROR mapDTO: Could not create placeholder domain.Currency %s: %v", cc.String(), ncErr)
				continue
			}
			domainCurrencies = append(domainCurrencies, *newCurrency)

		} else if err != nil {
			log.Printf("ERROR mapDTO: Failed to lookup currency %s for country %s: %v. Skipping this currency.", cc.String(), dto.CountryCode, err)
			continue

		} else if currency != nil { // meaning found in db
			domainCurrencies = append(domainCurrencies, *currency)
		}
	}

	return domain.NewCountry(countryCode, dto.CountryName, domainCurrencies, wikiID)
}
