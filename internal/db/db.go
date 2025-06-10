package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	//_ "github.com/godror/godror"    // Oracle driver
	_ "github.com/lib/pq"           // Postgres driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Config holds database configuration read from environment or file
type Config struct {
	Driver string // "sqlite", "postgres", or "oracle"
	DSN    string // Data Source Name
}

// LoadConfigPrefix reads e.g. SALES_DB_DRIVER & SALES_DB_DSN
func LoadConfigPrefix(prefix string) Config {
	drv := os.Getenv(prefix + "DB_DRIVER")
	if drv == "" {
		drv = "sqlite"
	}
	dsn := os.Getenv(prefix + "DB_DSN")
	if dsn == "" {
		// fallback file for sqlite
		dsn = fmt.Sprintf("data/%s.db?_foreign_keys=1", strings.ToLower(prefix))
	}
	return Config{Driver: drv, DSN: dsn}
}

func Open(cfg Config) (*sql.DB, error) {
	switch cfg.Driver {
	case "sqlite":
		return sql.Open("sqlite3", cfg.DSN)
	case "postgres":
		return sql.Open("postgres", cfg.DSN)
	case "oracle":
		return sql.Open("godror", cfg.DSN)
	default:
		return nil, fmt.Errorf("unsupported driver: %s", cfg.Driver)
	}
}
