package sqlxadapter_test

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"

	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// runs once at module load and changes working directory to project root, adjust the ../../.. as needed
func init() {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(filename, "../../../../../..")

	err := os.Chdir(dir)

	if err != nil {
		panic(err)
	}
}

func setupTestDB(t *testing.T) (*sqlx.DB, func()) {
	t.Helper()

	// dsn := "file::memory:?cache=shared"
	// tests manage the temp dir so it self cleans after testing
	// sqlite pragma options: fkeys on and wal on

	tempDBFilePath := filepath.Join(t.TempDir(), fmt.Sprintf("test_db_%d.sqlite", time.Now().UnixNano()))

	dbOptions := "?cache=shared&_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"

	dsn := fmt.Sprintf("file:%s%s", tempDBFilePath, dbOptions)

	tempDB, err := sql.Open("sqlite3", dsn)
	if err != nil {
		t.Fatalf("failed to open test db %v", err)
	}

	// anybody home?
	if err := tempDB.Ping(); err != nil {
		tempDB.Close()
		t.Fatalf("failed to ping test db")
	}

	migrationsPath := "./migrations" // relative to test file?

	// does migration files folder URL exist?
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		tempDB.Close()
		t.Fatalf("failed to find migrations %v", err)
	}

	migrationsURL := fmt.Sprintf("file://%s", migrationsPath)
	migrateDSN := fmt.Sprintf("sqlite3://%s", tempDBFilePath)

	// create new migration instance
	m, err := migrate.New(migrationsURL, migrateDSN)
	if err != nil {
		tempDB.Close()
		os.Remove(tempDBFilePath)
		t.Fatalf("failed to create migration instance: %v", err)
	}

	// apply up migrations
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		m.Close()
		tempDB.Close()
		os.Remove(tempDBFilePath)

		t.Fatalf("failed to apply migrations: %v", err)
	}

	log.Printf("Migrations applied to test db")

	// call sqlx
	sqlxDB := sqlx.NewDb(tempDB, "sqlite3")

	cleanup := func() {
		log.Println("Cleaning temp db")

		// reverse the migrations
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			m.Close()
		}

		sqlxDB.Close()
		os.Remove(tempDBFilePath)
	}

	return sqlxDB, cleanup
}

func TestCurrencyRepository(t *testing.T) {
	t.Run("Save and GetByID", func(t *testing.T) {
		_, cleanup := setupTestDB(t)
		defer cleanup()

		fmt.Println("got to the test")
	})

}
