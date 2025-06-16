package sqlitedbadapter

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Sentinel errors for unit testing later
var (
	ErrDBFolder        = errors.New("failed to create DB folder")
	ErrDBOpen          = errors.New("failed to open sqlite3 DB")
	ErrDBPing          = errors.New("failed to ping sqlite3 DB")
	ErrMigrationDriver = errors.New("failed to create migration driver")
	ErrMigrationRun    = errors.New("failed to run migrations")
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
	dsn := fmt.Sprintf("file:%s?cache=shared&_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)", dbFilePath)
	db, err := sql.Open("sqlite3", dsn)

	if err != nil {
		return nil, fmt.Errorf("%w %s: %w", ErrDBOpen, dbFilePath, err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%w %w", ErrDBPing, err)
	}

	log.Printf("Successfully connected to sqlite3 DB: %s", dbFilePath)
	return db, nil
}

// RunMigrations applies pending migrations
func RunMigrations(db *sql.DB, migrationsPath string) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrMigrationDriver, err)
	}

	sourceURL := "file://" + migrationsPath

	m, err := migrate.NewWithDatabaseInstance(sourceURL, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("%w: failed to create migration instance: %w", ErrMigrationDriver, err)
	}

	log.Printf("Running migrations from: %s", migrationsPath)

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%w: %w", ErrMigrationRun, err)
	}

	log.Println("Database migration check completed successfully.")
	return nil
}

// func CreateSchema(db *sql.DB) error {
// 	personSchema := `
// 		CREATE TABLE IF NOT EXISTS person (
// 			id BLOB(16) PRIMARY KEY,
// 			first_name VARCHAR(40) NOT NULL,
// 			last_name VARCHAR(40) NOT NULL,
// 			email VARCHAR(255) UNIQUE NOT NULL,
// 			dob DATETIME NOT NULL
// 		);`

// 	_, err := db.Exec(personSchema)
// 	if err != nil {
// 		return fmt.Errorf("%w: %w", ErrSchema, err)
// 	}

// 	log.Println("'person' schema created successfully")

// 	return nil
// }
