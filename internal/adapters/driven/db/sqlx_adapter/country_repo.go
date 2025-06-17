package sqlxadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
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

// Save writes a (domain object) Country into the DB returning the result of a DB get of the written instance
func (r *CountryRepo) Save(ctx context.Context, country *domain.Country) (*domain.Country, error) {

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: beginning transaction for saving country %s: %v", dbcommon.ErrTransactionBegin, country.Code(), err)
	}

	// defer a rollback
	defer func() {
		// if there was a panic
		if p := recover(); p != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("ERROR: transaction rollback failed for country %s after error %v: %v", country.Code(), err, rbErr)
			}
			panic(p) // repanic anyway
		}
		// if an error occurred, rollback
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("ERROR: transaction rollback failed for country %s after error %v: %v", country.Code(), err, rbErr)
			}
		}
	}()

	countryModel := toModelCountry(country)
	if countryModel == nil {
		err = dbcommon.ErrConvertNilCountry
		return nil, err
	}

	saveCountryQuery, err := GetQuery("SaveCountry")
	if err != nil {
		err = fmt.Errorf("SaveCurrency query retrieval failed: %w", err)
		return nil, err
	}

	_, err = tx.NamedExecContext(ctx, saveCountryQuery, countryModel)
	if err != nil {
		err = fmt.Errorf("%w for country code %s: %s: %v", dbcommon.ErrSQLxSaveCountry, country.Code(), country.Name(), err)
		return nil, err
	}

	// delete previous country associations and recreate
	deleteCountryCurrencyJoinsQuery, err := GetQuery("DeleteCountryCurrencyJoins")
	if err != nil {
		err = fmt.Errorf("DeleteCountryCurrencyJoins query retrieval failed: %w", err)
		return nil, err
	}

	_, err = tx.NamedExecContext(ctx, deleteCountryCurrencyJoinsQuery, countryModel)
	if err != nil {
		err = fmt.Errorf("%w for country code %s: %v", dbcommon.ErrSQLxDeleteJoins, country.Name(), err)
		return nil, err
	}

	if len(country.Currencies()) > 0 {

		saveCountryCurrencyPairQuery, err := GetQuery("SaveCountryCurrencyPair")
		if err != nil {
			return nil, fmt.Errorf("SaveCountryCurrencyPair query retrieval failed: %w", err)
		}

		for _, c := range country.Currencies() {
			row := struct {
				CountryCode  domain.CountryCode  `db:"country_code"`
				CurrencyCode domain.CurrencyCode `db:"currency_code"`
			}{
				CountryCode:  country.Code(),
				CurrencyCode: c.Code(),
			}
			_, err = tx.NamedExecContext(ctx, saveCountryCurrencyPairQuery, row)

			if err != nil {
				err = fmt.Errorf("%w for country/currency pair %s: %s: %v", dbcommon.ErrSQLxSaveCountryCurrency, country.Code(), c.Code().String(), err)
				return nil, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		// if commit fails a rollback will be attempted from defer
		return nil, fmt.Errorf("%w: committing transaction for country %s: %v", dbcommon.ErrTransactionCommit, country.Code(), err)
	}

	createdCountry, err := r.GetByID(ctx, country.Code())
	if err != nil {
		return nil, fmt.Errorf("%w for country %s: %v", dbcommon.ErrSQLxSavedButNotInDB, country.Name(), err)
	}

	log.Printf("Country %s and its currencies saved/updated successfully", country.Code())
	return createdCountry, nil
}

// GetByID returns a Country instance if its country code exists
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

	retrievedCountry, err := countryModel.toDomainCountry(cModels)
	if err != nil {
		return nil, fmt.Errorf("%w for country %s: %v", dbcommon.ErrConvertToCountry, ccStr, err)
	}

	return retrievedCountry, nil
}

// CountAll returns the number of entries in the Country table
func (r *CountryRepo) CountAll(ctx context.Context) (int, error) {
	countAllCountriesQuery, err := GetQuery("CountAllCountries")
	if err != nil {
		return 0, fmt.Errorf("CountAllCountries query retrieval: %w", err)
	}

	var count int
	err = r.db.GetContext(ctx, &count, countAllCountriesQuery)

	if err != nil {
		return 0, fmt.Errorf("%w: counting all countries: %v", dbcommon.ErrSQLxQueryFailed, err)
	}

	return count, nil
}

// GetRandom returns a random Country from DB
func (r *CountryRepo) GetRandom(ctx context.Context) (*domain.Country, error) {
	getRandomCountryQuery, err := GetQuery("GetRandomCountry")
	if err != nil {
		return nil, fmt.Errorf("GetRandomCountry query retrieval: %w", err)
	}

	var cModel CountryModel // not a pointer
	err = r.db.GetContext(ctx, &cModel, getRandomCountryQuery)
	if err != nil {
		return nil, fmt.Errorf("%w: getting random country: %v", dbcommon.ErrSQLxQueryFailed, err)
	}

	currenciesQuery, err := GetQuery("GetCurrenciesForCountry")
	if err != nil {
		return nil, fmt.Errorf("GetCurrenciesForCountry query retrieval failed: %w", err)
	}

	ccStr := cModel.Code.String()
	var cModels []CurrencyModel
	// SelectContext needs a pointer to a slice...
	err = r.db.SelectContext(ctx, &cModels, currenciesQuery, ccStr)

	if err != nil {
		return nil, fmt.Errorf("%w: fetching currencies for country '%s': %v", dbcommon.ErrSQLxQueryFailed, ccStr, err)
	}

	// create the domain Country object
	result, err := cModel.toDomainCountry(cModels)
	if err != nil {
		return nil, fmt.Errorf("%w for country %s: %v", dbcommon.ErrConvertToCountry, ccStr, err)
	}

	return result, nil
}
