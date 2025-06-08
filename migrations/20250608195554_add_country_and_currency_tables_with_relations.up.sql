CREATE TABLE IF NOT EXISTS country (
    code CHAR(2) PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    wikidataid VARCHAR(10) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS currency(
    code CHAR(3) PRIMARY KEY,
    name VARCHAR(50) UNIQUE
);

CREATE TABLE IF NOT EXISTS country_currency(
    country_code CHAR(2) NOT NULL,
    currency_code CHAR(3) NOT NULL,
    CONSTRAINT pk_country_currency PRIMARY KEY (country_code, currency_code),
    CONSTRAINT fk_country_currency_country FOREIGN KEY (country_code) REFERENCES country (code) ON DELETE RESTRICT,
    CONSTRAINT fk_country_currency_currency FOREIGN KEY (currency_code) REFERENCES currency (code) ON DELETE RESTRICT
);

-- BEGIN TRANSACTION;

ALTER TABLE person RENAME TO person_old;

CREATE TABLE IF NOT EXISTS person (
    id BLOB(16) PRIMARY KEY,
    first_name VARCHAR(40) NOT NULL CHECK(LENGTH(first_name) <= 40),
    last_name VARCHAR(40) NOT NULL CHECK(LENGTH(last_name) <= 40),
    email VARCHAR(255) UNIQUE NOT NULL CHECK(LENGTH(email) <= 255),
    dob DATETIME NOT NULL
        -- IMPORTANT: Use strftime to force UTC default
        CHECK (datetime(dob) IS NOT NULL AND substr(dob, -1) = 'Z'),
    created_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
        -- Also UTC'
        CHECK (datetime(created_at) IS NOT NULL AND substr(created_at, -1) = 'Z'),
    updated_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%SZ', 'now'))
        -- Also enforce UTC
        CHECK (datetime(updated_at) IS NOT NULL AND substr(updated_at, -1) = 'Z'),

    -- new stuff here
    country_code CHAR(2),
    CONSTRAINT fk_person_country FOREIGN KEY (country_code) REFERENCES country (code) ON DELETE RESTRICT
);

INSERT INTO person (id, first_name, last_name, email, dob, created_at, updated_at) SELECT id, first_name, last_name, email, dob, created_at, updated_at FROM person_old;

DROP TABLE person_old;
-- COMMIT;

DROP TRIGGER IF EXISTS person_updated_at;

-- automate updated_at
CREATE TRIGGER person_updated_at
AFTER UPDATE ON person
FOR EACH ROW
BEGIN
    UPDATE person SET updated_at = strftime('%Y-%m-%dT%H:%M:%SZ', 'now') WHERE id = OLD.id;
END;
