-- name: SaveCurrency
-- Inserts a new currency or updates an existing one if the code matches.
INSERT INTO currency (code, name)
VALUES (:code, :name)
ON CONFLICT(code) DO UPDATE SET
    name = excluded.name;

-- name: GetCurrencyByCode
-- Selects a currency by its unique code.
SELECT code, name
FROM currency
WHERE code = ?;

-- name: CountAllCurrencies
-- Counts all currencies in the table.
SELECT COUNT(*) FROM currency;

-- name: GetRandomCurrency
-- Selects a random currency from the table (SQLite specific).
SELECT code, name FROM currency LIMIT 1 OFFSET ?

-- name: ListAllCurrencies
-- Selects all currencies.
SELECT code, name FROM currency ORDER BY name;
