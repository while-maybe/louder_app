package sqlxadapter_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	sqlxadapter "louder/internal/adapters/driven/db/sqlx_adapter"
	"louder/internal/core/domain"
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

// runs once at module load and changes working directory to project root, adjust the ../../../../.. as needed
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

	tempDBFilePath := filepath.Join(t.TempDir(), fmt.Sprintf("test_db_%d.sqlite", time.Now().UnixNano()))
	// sqlite pragma options: fkeys on and wal on
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

	migrationsPath := "./migrations" // bless init() on top

	// does migration files folder URL exist?
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		tempDB.Close()
		t.Fatalf("failed to find migrations %v", err)
	}

	migrationsURL := fmt.Sprintf("file://%s", migrationsPath)
	migrateDSN := fmt.Sprintf("sqlite3://%s", tempDBFilePath)

	// create new migration instance
	// the sqlite3 driver is NOT mattn's - migrate has its own :/
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

	log.Printf("Migrations up on test db")

	// call sqlx
	sqlxDB := sqlx.NewDb(tempDB, "sqlite3")

	cleanupTestDB := func() {
		log.Println("Cleaning up test db")

		// reverse the migrations
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			m.Close()
		}

		log.Printf("Migrations down on test db")

		sqlxDB.Close()
		os.Remove(tempDBFilePath)
	}

	return sqlxDB, cleanupTestDB
}

func TestSaveAndGetByID(t *testing.T) {
	eurCode, _ := domain.NewCurrencyCode("eur")

	tt := map[string]struct {
		currToSave   *domain.Currency
		wantCurrID   domain.CurrencyCode
		wantCurrName string

		wantSaveErr error
		wantGetErr  error
	}{
		"Save/Get new currency ok": {
			currToSave: func() *domain.Currency {
				c, _ := domain.NewCurrency(eurCode, "Euro")
				return c
			}(), // call it here
			wantCurrID:   "EUR",
			wantCurrName: "Euro",
			wantSaveErr:  nil,
			wantGetErr:   nil,
		},
		// other cases, check coverage?
	}

	// iterate through tests here
	for name, tc := range tt {

		t.Run(name, func(t *testing.T) {
			// for db ops, best practice is to clean db and recreate per test...
			db, cleanup := setupTestDB(t)
			defer cleanup()

			// NewCurrencyRepo needs a *sql.DB which is encapsulated inside *sqlx.db.DB, hence db.DB
			repo, err := sqlxadapter.NewCurrencyRepo(db.DB)
			if err != nil {
				t.Fatalf("failed to create currency repo: %v", err)
			}

			ctx := context.Background()

			// testing starts here
			gotCurr, saveErr := repo.Save(ctx, tc.currToSave)

			// save errors match
			if tc.wantSaveErr != saveErr {
				t.Fatalf("unexpected error: expected %v got %v", tc.wantSaveErr, saveErr)
			}

			// currency code matches
			if tc.wantCurrID != gotCurr.Code() {
				t.Fatalf("error unexpected currency code: expected %v got %v", tc.wantCurrID, gotCurr.Code())
			}

			// currency name matches
			if tc.wantCurrName != gotCurr.Name() {
				t.Fatalf("error unexpected currency name: expected %v got %v", tc.wantCurrName, gotCurr.Name())
			}
		})
	}
}
