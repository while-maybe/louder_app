package dbcommon

import "errors"

// errors for person
var (
	ErrHydrateWithNil   = errors.New("error attempted to hydrate person without data")
	ErrSaveNilPerson    = errors.New("error cannot save nil person to DB model")
	ErrConvertNilPerson = errors.New("error convert nil person to DB model")
	ErrEmptyID          = errors.New("error given ID is empty")
	ErrInvalidID        = errors.New("error invalid person ID format")
	ErrNilDomainPerson  = errors.New("error conversion returned nil domain person without error")
	ErrConvertToPerson  = errors.New("error converting SQLx/Bun data to a person")
)

// common db errors
var (
	ErrNotFound             = errors.New("error cannot find this ID in DB")
	ErrSavedButNotInDB      = errors.New("error SQLx/Bun person saved but can't find in DB")
	ErrDBQueryFailed        = errors.New("error query has failed")
	ErrSQLxSavedButNotInDB  = errors.New("error SQLx entity saved but could not get from DB")
	ErrDomainCreation       = errors.New("error cannot create domain entity")
	ErrTransactionBegin     = errors.New("error starting SQLx transaction")
	ErrTransactionCommit    = errors.New("error committing SQLx transaction")
	ErrSQLxQueryFailed      = errors.New("error SQLx failed to run the query")
	ErrSQLxNoRowsAffected   = errors.New("error SQLx could not get rows affected")
	ErrSQLxZeroRowsAffected = errors.New("error SQLx got 0 rows affected. Upsert?")
)

// errors for Country
var (
	ErrConvertCountry    = errors.New("error converting SQLx/Bun data to a country")
	ErrConvertNilCountry = errors.New("error converting nil country to DB model")
	ErrSQLxSaveCountry   = errors.New("error could not save country to DB model")

	ErrSQLxSaveCountryCurrency = errors.New("error could not save country/currency pair to DB model")
	ErrNoCountryCode           = errors.New("error currency code must be provided")
	ErrConvertToCountry        = errors.New("error converting DB data to country model")
	ErrSQLxDeleteJoins         = errors.New("error deleting existing country/currency entries in DB")
)

// errors for Currency
var (
	ErrConvertCurrency    = errors.New("error converting SQLx/Bun data to a currency")
	ErrConvertNilCurrency = errors.New("error converting nil currency to DB model")
	ErrSaveCurrency       = errors.New("error could not save currency to DB model")
	ErrNoCurrencyCode     = errors.New("error currency code must be provided")
	ErrSQLxNotFound       = errors.New("error SQLx value not in DB")
	ErrConvertToCurrency  = errors.New("error converting DB data to currency model")
)
