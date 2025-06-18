package geodbclient

import (
	"context"
	"errors"
	"fmt"
	"log"
	"louder/internal/adapters/driven/db/dbcommon"
	"louder/internal/core/domain"
)

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
