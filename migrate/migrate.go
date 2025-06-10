package migrate

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// MigrateSQLite will run every .sql file in dir (in alphabetical order).
func MigrateSQLite(db *sql.DB, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading migrations dir: %w", err)
	}

	// Collect .sql files
	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	// Apply each
	for _, fname := range files {
		path := filepath.Join(dir, fname)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", fname, err)
		}
		stmts := strings.Split(string(content), ";")
		for _, stmt := range stmts {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" || strings.HasPrefix(stmt, "--") {
				continue
			}
			if _, err := db.Exec(stmt); err != nil {
				return fmt.Errorf("%s: exec %q: %w", fname, stmt, err)
			}
		}
	}
	return nil
}
