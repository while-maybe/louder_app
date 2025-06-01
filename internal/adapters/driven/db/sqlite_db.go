package sqlitedbadapter

import (
	"database/sql"
	"fmt"
	"log"
	dbcommon "louder/internal/adapters/driven/db/db_common"
	"os"
	"path/filepath"
)

// Sentinel errors for unit testing later
const (
	ErrDBFolder = dbcommon.Error("failed to create DB folder")
	ErrDBOpen   = dbcommon.Error("failed to open sqlite3 DB")
	ErrDBPing   = dbcommon.Error("failed to ping sqlite3 DB")
	ErrSchema   = dbcommon.Error("failed to create persons schema")
)

func Init(dbFilePath string) (*sql.DB, error) {
	dbDir := filepath.Dir(dbFilePath)

	if _, err := os.Stat(dbDir); os.IsNotExist(err) {

		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("%w %s: %w", ErrDBFolder, dbDir, err)
		}
	}
	log.Printf("DB directory '%s' created", dbDir)

	// sqlite pragma options: fkeys on and wal on
	dns := fmt.Sprintf("file:%s?_foreign_keys=on&_journal_model=WAL", dbFilePath)
	db, err := sql.Open("sqlite3", dns)

	if err != nil {
		return nil, fmt.Errorf("%w %s: %w", ErrDBOpen, dbFilePath, err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%w %w", ErrDBPing, err)
	}

	log.Printf("Successfully connected to sqlite3 DB: %s", dbFilePath)
	return db, nil
}

func CreateSchema(db *sql.DB) error {
	personSchema := `
		CREATE TABLE IF NOT EXISTS person (
			id VARCHAR(16) PRIMARY KEY,
			first_name VARCHAR(40) NOT NULL,
			last_name VARCHAR(40) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			dob VARCHAR(30) NOT NULL
		);`

	_, err := db.Exec(personSchema)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSchema, err)
	}

	log.Println("'person' schema created successfully")

	return nil
}
