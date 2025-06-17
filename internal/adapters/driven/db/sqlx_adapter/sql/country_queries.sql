-- name: SaveCountry
-- Inserts a new country or updates an existing one if the code matches.
INSERT INTO country (code, name, wikidataid)
VALUES (:code, :name, :wikidataid)
ON CONFLICT(code) DO UPDATE SET
    name = excluded.name,
    wikidataid = excluded.wikidataid;

-- name: GetCountryByCode
-- Returns a Country given its 3 digit ISO code
SELECT code, name, wikidataid FROM country WHERE code = ?;

-- name: GetCurrenciesForCountry
-- Returns all currencies for a given country code
SELECT currency_code FROM country_currency WHERE country_code = ?;

-- name: CountAllCountries
-- Return the count of all existing countries
SELECT COUNT(*) FROM country;

-- name GetRandomCountry
-- Returns one country at random
SELECT code, name, wikidataid FROM country ORDER BY RANDOM() LIMIT 1;

-- name ListAllCountries
-- Returns all countries
SELECT code, name, wikidataid FROM country ORDER BY name;

-- name: SaveCountryCurrencyPair
-- Inserts a new country/currency pair or updates an existing one if the code matches.
INSERT INTO country_currency (country_code, currency_code)
VALUES (:country_code, :currency_code)
ON CONFLICT (country_code, currency_code) DO NOTHING;

-- name DeleteCountryCurrencyJoins
-- Deletes rows representing all currencies associated with a Country
DELETE FROM country_currencies WHERE country_code = :country_code;