package sqlxadapter

import (
	"fmt"
	dbcommon "louder/internal/adapters/driven/db/dbcommon"
	"louder/internal/core/domain"
)

type CountryModel struct {
	Code       domain.CountryCode `db:"country_code"`
	Name       string             `db:"name"`
	WikiDataID domain.WikiCode    `db:"wikidataid"`
}

// toModelCountry takes a Country domain entity and returns its equivalent SQLx model
func toModelCountry(c *domain.Country) *CountryModel {
	if c == nil {
		return nil
	}

	return &CountryModel{
		Code:       c.Code(),
		Name:       c.Name(),
		WikiDataID: c.WikiId(),
	}
}

// toDomainCountry takes a SQLx country model and returns its equivalent domain entity
func (m *CountryModel) toDomainCountry(currencyModels []CurrencyModel) (*domain.Country, error) {
	if m == nil {
		return nil, dbcommon.ErrConvertCountry
	}

	domainCurrencies := make([]domain.Currency, 0, len(currencyModels))

	for _, cModel := range currencyModels {
		tempDomainCurrency, err := (&cModel).toDomainCurrency()

		if err != nil {
			// Return an error, providing context about which currency conversion failed.
			return nil, fmt.Errorf("failed to convert currency model (code: %s) to domain currency: %w",
				cModel.Code.String(), err)
		}

		if tempDomainCurrency == nil {
			return nil, fmt.Errorf("internal error: toDomainCurrency for model (code: %s) returned nil domain currency without error", cModel.Code.String())
		}

		domainCurrencies = append(domainCurrencies, *tempDomainCurrency)
	}

	createdCountry, err := domain.NewCountry(m.Code, m.Name, domainCurrencies, m.WikiDataID)
	if err != nil {
		return nil, fmt.Errorf("%w: while creating domain country from model (code: %s): %v", dbcommon.ErrDomainCreation, m.Code.String(), err)
	}

	return createdCountry, nil
}
