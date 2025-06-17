package sqlxadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"louder/internal/adapters/driven/db/dbcommon"
	"louder/internal/core/domain"
	"louder/internal/core/service/countrycore"

	"github.com/jmoiron/sqlx"
)

type CountryRepo struct {
	db *sqlx.DB
}

// we attempt to assign a value of type pointer *CountryRepo (which we assert to be nil) to a var of the interface type to make sure the value implements it
// ensure CountryRepo implements the Port (safety check)
var _ countrycore.Repository = (*CountryRepo)(nil)

// Factory function (generator) for an interface of the type CountryRepo - return an interface here, not a instance of CountryRepo
func NewCountryRepo(sqldb *sql.DB) (countrycore.Repository, error) {
	db := sqlx.NewDb(sqldb, "sqlite3")
	return &CountryRepo{db: db}, nil
}

func (r *CountryRepo) Save(ctx context.Context, country *domain.Country) (*domain.Country, error) {
	countryModel := toModelCountry(country)
	if countryModel == nil {
		return nil, dbcommon.ErrConvertNilCountry
	}

	query, err := GetQuery("SaveCountry")
	if err != nil {
		return nil, fmt.Errorf("SaveCurrency query retrieval: %w", err)
	}

	_, err = r.db.NamedExecContext(ctx, query, countryModel)
	if err != nil {
		return nil, fmt.Errorf("%w for country code %s: %s: %v", dbcommon.ErrSQLxSaveCountry, country.Code(), country.Name(), err)
	}

	query, err = GetQuery("SaveCountryCurrencyPair")
	if err != nil {
		return nil, fmt.Errorf("SaveCountryCurrencyPair query retrieval failed: %w", err)
	}

	for _, c := range country.Currencies() {
		row := struct {
			CountryCode  domain.CountryCode  `db:"country_code"`
			CurrencyCode domain.CurrencyCode `db:"currency_code"`
		}{
			CountryCode:  country.Code(),
			CurrencyCode: c,
		}
		_, err := r.db.NamedExecContext(ctx, query, row)

		if err != nil {
			return nil, fmt.Errorf("%w for country/currency pair %s: %s: %v", dbcommon.ErrSQLxSaveCountryCurrency, country.Code(), c.String(), err)
		}

	}

	query, err = GetQuery("GetCountryByCode")
	if err != nil {
		return nil, fmt.Errorf("GetCountryByCode query retrieval failed: %w", err)
	}

	createdCountry, err := r.GetByID()

}

func (r *CountryRepo) GetByID(ctx context.Context, cc domain.CountryCode) (*domain.Country, error) {
	ccStr := cc.String()
	if ccStr == "" {
		return nil, dbcommon.ErrNoCountryCode
	}

	countryQuery, err := GetQuery("GetCountryByCode")
	if err != nil {
		return nil, fmt.Errorf("GetCountryByCode query retrieval failed: %w", err)
	}
	currenciesQuery, err := GetQuery("GetCurrenciesForCountry")
	if err != nil {
		return nil, fmt.Errorf("GetCurrenciesForCountry query retrieval failed: %w", err)
	}

	var countryModel CountryModel
	err = r.db.GetContext(ctx, &countryModel, countryQuery, ccStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// no results
			return nil, fmt.Errorf("%w for country code %s: %v", dbcommon.ErrSQLxNotFound, ccStr, err)
		}
		// query fails
		return nil, fmt.Errorf("%w: %v", dbcommon.ErrSQLxQueryFailed, err)
	}

	var cModels []CurrencyModel
	// SelectContext needs a pointer to a slice...
	err = r.db.SelectContext(ctx, &cModels, currenciesQuery, ccStr)

	if err != nil {
		return nil, fmt.Errorf("%w: fetching currencies for country '%s': %v", dbcommon.ErrSQLxQueryFailed, ccStr, err)
	}

	if cModels == nil {
		cModels = []CurrencyModel{}
	}

	retrievedCountry, err := countryModel.toDomainCountry(cModels)
	if err != nil {
		return nil, fmt.Errorf("%w for country %s: %v", dbcommon.ErrConvertToCountry, ccStr, err)
	}

	return retrievedCountry, nil
}

// Save(ctx context.Context, country *domain.Country) (*domain.Country, error)
// GetByID(ctx context.Context, cc domain.CountryCode) (*domain.Country, error) // ID is the Country's ISO code
// CountAll(ctx context.Context) (uint, error)
// GetRandom(ctx context.Context) (*domain.Country, error)
