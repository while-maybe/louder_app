package sqlxadapter

import (
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
)

//go:embed sql/*.sql
var queryFS embed.FS

var (
	queries map[string]string
	once    sync.Once
	initErr error
)

// initQueries loads and parses all SQL queries from the embedded filesystem. It's called once using sync.Once to ensure thread-safety and avoid re-parsing.
func initQueries() {
	once.Do(func() {
		queries = make(map[string]string)

		err := fs.WalkDir(queryFS, "sql", func(path string, d fs.DirEntry, walkError error) error {
			if walkError != nil {
				return fmt.Errorf("error walking directory %s: %w", path, walkError)
			}

			if !d.IsDir() && filepath.Ext(path) == ".sql" {
				content, readErr := queryFS.ReadFile(path)

				if readErr != nil {
					return fmt.Errorf("failed to read embedded sql file %s: %w", path, readErr)
				}

				parseContentIntoQueries(string(content))
			}

			return nil
		})

		if err != nil {
			initErr = fmt.Errorf("failed to initialize SQL queries for sqlx_adapter: %w", err)
			// yes, panic baby
			panic(initErr)
		}
	})
}

// parseContentIntoQueries parses a single SQL file's content and populates the queries map.
func parseContentIntoQueries(content string) {
	// statements := strings.Split(content, "--name:")
	// Go 1.24 implements SplitSeq which does not create a slice for the results, instead, returns an interator which doesn't allocate.

	statements := strings.SplitSeq(content, "-- name:")

	for s := range statements {
		s = strings.TrimSpace(s)

		if s == "" {
			continue
		}
		parts := strings.SplitN(s, "\n", 2)
		if len(parts) < 2 {
			// Log or handle malformed query name line
			fmt.Printf("Warning: Malformed query block in SQL content (missing newline after -- name:): %s\n", s)
			continue
		}

		nameParts := strings.Fields(parts[0]) // Get "Name" and potentially description
		if len(nameParts) == 0 {
			fmt.Printf("Warning: Malformed query block in SQL content (empty -- name: line): %s\n", parts[0])
			continue
		}

		queryName := nameParts[0]
		querySQL := strings.TrimSpace(parts[1])

		if _, exists := queries[queryName]; exists {
			// Handle duplicate query names, e.g., panic or log a warning. This indicates a potential issue in .sql files.
			initErr = fmt.Errorf("duplicate SQL query name detected: %s. Check your .sql files", queryName)

			// panic here too
			dupError := fmt.Sprintf("ERROR: Duplicate SQL query name detected: %s. Check your .sql files.\n", queryName)
			panic(dupError)
		}

		queries[queryName] = querySQL
	}
}

// GetQuery retrieves a named query. It ensures queries are initialized.
func GetQuery(name string) (string, error) {
	initQueries()

	if initErr != nil {
		return "", fmt.Errorf("failed to initialize queries: %w (query name: %s)", initErr, name)
	}

	query, ok := queries[name]
	if !ok {
		return "", fmt.Errorf("sql query '%s' not found in sqlx_adapter", name)
	}

	return query, nil
}
