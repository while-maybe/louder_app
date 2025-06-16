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
	"math/rand/v2"

	"github.com/jmoiron/sqlx"
)

type CurrencyRepo struct {
	db *sqlx.DB
}

// ensure SQLxCurrencyRepo implements the Port (safety check)
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

	// write the query
	query := `
		INSERT INTO currency (code, name)
		VALUES (:code, :name)
		ON CONFLICT(code) DO UPDATE SET
			name = excluded.name;`

	// run the query and get the result (and check for errors)
	_, err := r.db.NamedExecContext(ctx, query, sqlxModel)
	if err != nil {
		return nil, fmt.Errorf("%w for currency code %s: %s", dbcommon.ErrSQLxSaveCurrency, currency.Code(), currency.Name())
	}

	// checking for rows affected == 0 might not be great here as the query has a DO UPDATE SET, so if an upsert occurs, rowsaffected would have returned 0 and that is not an erorr

	var createdCurrency *domain.Currency

	createdCurrency, err = r.GetByID(ctx, currency.Code())
	if err != nil {
		return nil, fmt.Errorf("%w for currency code %s: %w", dbcommon.ErrSQLxSavedButNotInDB, currency.Code(), err)
	}
	log.Printf("Currency %s inserted/updated with success", currency.Code())

	return createdCurrency, nil
}

func (r *CurrencyRepo) GetByID(ctx context.Context, cc domain.CurrencyCode) (*domain.Currency, error) {
	givenCode := string(cc)
	if string(givenCode) == "" {
		return nil, dbcommon.ErrNoCurrencyCode
	}

	query := `
		SELECT code, name
		FROM currency
		WHERE code = ?;
	`

	var sqlxModel CurrencyModel

	err := r.db.GetContext(ctx, &sqlxModel, query, givenCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// no results
			return nil, fmt.Errorf("%w for currency code %s: %w", dbcommon.ErrSQLxNotFound, givenCode, err)
		}
		// query fails
		return nil, fmt.Errorf("%w: %w", dbcommon.ErrSQLxQueryFailed, err)
	}

	retrievedCurrency, err := sqlxModel.toDomainCurrency()
	if err != nil {
		return nil, fmt.Errorf("%w for currency %s: %w", dbcommon.ErrConvertToCurrency, givenCode, err)
	}
	return retrievedCurrency, nil
}

func (r *CurrencyRepo) CountAll(ctx context.Context) (int, error) {
	// the count will include entries with NULL, if an issue use SELECT(code)
	query := `SELECT COUNT(*) FROM currency`

	var count int
	err := r.db.GetContext(ctx, &count, query)

	if err != nil {
		return 0, fmt.Errorf("%w: counting all currencies: %w", dbcommon.ErrSQLxQueryFailed, err)
	}

	return count, nil
}

func (r *CurrencyRepo) GetRandom(ctx context.Context) (*domain.Currency, error) {
	// I don't like the idea of SELECT code, name FROM currency ORDER BY RANDOM() LIMIT 1;
	// Postgres would have random selection so considering this is for sqlite, I'll accept it

	currenciesCount, err := r.CountAll(ctx)
	if err != nil {
		return nil, err
	}
	if currenciesCount == 0 {
		return nil, fmt.Errorf("%w: no currencies available in the database", dbcommon.ErrSQLxNotFound)
	}

	randomOffset := rand.IntN(currenciesCount)

	// Get a row at a specific offset
	query := `SELECT code, name FROM currency LIMIT 1 OFFSET ?`
	queryResult := r.db.QueryRowContext(ctx, query, randomOffset)

	var row *CurrencyModel

	err = queryResult.Scan(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: %w", dbcommon.ErrSQLxNotFound, err)
		}
		return nil, fmt.Errorf("%w: %w", dbcommon.ErrSQLxQueryFailed, err)
	}

	// convert to a domain Currency
	result, err := row.toDomainCurrency()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", dbcommon.ErrConvertToCurrency, err)
	}

	return result, nil
}

// var sqlxModel SQLxModelCurrency

// type Repository interface {

// 	// GetByName(ctx context.Context, name string) (*domain.Currency, error)
// 	// Search(ctx context.Context, terms string) ([]*domain.Currency, error) // get a list of countries when search terms are given, like Google?
// }
