package sqlxadapter

import (
	dbcommon "louder/internal/adapters/driven/db/dbcommon"
	"louder/internal/core/domain"
)

type SQLxModelCountry struct {
	Code       domain.CountryCode `db:"country_code"`
	Name       string             `db:"name"`
	WikiDataID domain.WikiCode    `db:"wikidataid"`
}

// toSQLxModelCountry takes a Country domain entity and returns its equivalent SQLx model
func toSQLxModelCountry(c *domain.Country) *SQLxModelCountry {
	if c == nil {
		return nil
	}

	return &SQLxModelCountry{
		Code:       c.Code(),
		Name:       c.Name(),
		WikiDataID: c.WikiId(),
	}
}

// toDomainCountry takes a SQLx country model and returns its equivalent domain entity
func (c *SQLxModelCountry) toDomainCountry() (*domain.Country, error) {
	if c == nil {
		return nil, dbcommon.ErrConvertCountry
	}

	return domain.NewCountry(c.Code, c.Name, []domain.CurrencyCode{}, c.WikiDataID), nil

}
