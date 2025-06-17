package sqlxadapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	dbcommon "louder/internal/adapters/driven/db/dbcommon"
	"louder/internal/core/domain"
	"louder/internal/core/service/currencycore"

	"github.com/jmoiron/sqlx"
)

type CurrencyRepo struct {
	db *sqlx.DB
}

// ensure CurrencyRepo implements the Port (safety check)
var _ currencycore.Repository = (*CurrencyRepo)(nil)

// return an interface here, not a instance of CurrencyRepo
func NewCurrencyRepo(sqldb *sql.DB) (currencycore.Repository, error) {
	db := sqlx.NewDb(sqldb, "sqlite3")
	return &CurrencyRepo{db: db}, nil
}

// Save places
func (r *CurrencyRepo) Save(ctx context.Context, currency *domain.Currency) (*domain.Currency, error) {
	// convert from domain model to sqlx model
	sqlxModel := toModelCurrency(currency)

	// check if it's nil
	if sqlxModel == nil {
		return nil, dbcommon.ErrConvertNilCurrency
	}

	query, err := GetQuery("SaveCurrency")
	if err != nil {
		return nil, fmt.Errorf("SaveCurrency query retrieval: %w", err)
	}

	// run the query and get the result (and check for errors)
	_, err = r.db.NamedExecContext(ctx, query, sqlxModel)
	if err != nil {
		return nil, fmt.Errorf("%w for currency code %s: %s: %v", dbcommon.ErrSQLxSaveCurrency, currency.Code(), currency.Name(), err)
	}

	// checking for rows affected == 0 might not be great here as the query has a DO UPDATE SET, so if an upsert occurs, rowsaffected would have returned 0 and that is not an erorr

	var createdCurrency *domain.Currency

	createdCurrency, err = r.GetByID(ctx, currency.Code())
	if err != nil {
		return nil, fmt.Errorf("%w for currency code %s: %v", dbcommon.ErrSQLxSavedButNotInDB, currency.Code(), err)
	}
	log.Printf("Currency %s inserted/updated with success", currency.Code())

	return createdCurrency, nil
}

func (r *CurrencyRepo) GetByID(ctx context.Context, cc domain.CurrencyCode) (*domain.Currency, error) {
	givenCode := string(cc)
	if string(givenCode) == "" {
		return nil, dbcommon.ErrNoCurrencyCode
	}

	query, err := GetQuery("GetCurrencyByCode")
	if err != nil {
		return nil, fmt.Errorf("%w: querying currency by code %s: %v", dbcommon.ErrSQLxQueryFailed, givenCode, err)
	}

	var sqlxModel CurrencyModel

	err = r.db.GetContext(ctx, &sqlxModel, query, givenCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// no results
			return nil, fmt.Errorf("%w for currency code %s: %v", dbcommon.ErrSQLxNotFound, givenCode, err)
		}
		// query fails
		return nil, fmt.Errorf("%w: %v", dbcommon.ErrSQLxQueryFailed, err)
	}

	retrievedCurrency, err := sqlxModel.toDomainCurrency()
	if err != nil {
		return nil, fmt.Errorf("%w for currency %s: %v", dbcommon.ErrConvertToCurrency, givenCode, err)
	}
	return retrievedCurrency, nil
}

func (r *CurrencyRepo) CountAll(ctx context.Context) (int, error) {
	// the count will include entries with NULL, if an issue use SELECT(code)

	query, err := GetQuery("CountAllCurrencies")
	if err != nil {
		return 0, fmt.Errorf("CountAllCurrencies query retrieval: %w", err)
	}

	var count int
	err = r.db.GetContext(ctx, &count, query)

	if err != nil {
		return 0, fmt.Errorf("%w: counting all currencies: %v", dbcommon.ErrSQLxQueryFailed, err)
	}

	return count, nil
}

func (r *CurrencyRepo) GetRandom(ctx context.Context) (*domain.Currency, error) {
	query, err := GetQuery("GetRandomCurrency")
	if err != nil {
		return nil, fmt.Errorf("GetRandomCurrency query retrieval: %w", err)
	}

	var row CurrencyModel // not a pointer
	err = r.db.GetContext(ctx, &row, query)
	if err != nil {
		return nil, fmt.Errorf("%w: getting random currency: %v", dbcommon.ErrSQLxQueryFailed, err)
	}

	// convert to a domain Currency
	result, err := row.toDomainCurrency()
	if err != nil {
		return nil, fmt.Errorf("%w: converting random currency: %v", dbcommon.ErrConvertToCurrency, err)
	}

	return result, nil
}

// TODO implement remaining port methods
// GetByName(ctx context.Context, name string) (*domain.Currency, error)

// get a list of countries when search terms are given, like Google?
// Search(ctx context.Context, terms string) ([]*domain.Currency, error)
