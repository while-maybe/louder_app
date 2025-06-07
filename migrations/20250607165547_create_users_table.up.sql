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
        CHECK (datetime(updated_at) IS NOT NULL AND substr(updated_at, -1) = 'Z')
);

-- automate updated_at
CREATE TRIGGER IF NOT EXISTS person_updated_at
AFTER UPDATE ON person
FOR EACH ROW
BEGIN
    UPDATE person SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;
