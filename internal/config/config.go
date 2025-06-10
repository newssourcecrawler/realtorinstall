// api/internal/config/config.go
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	// Per-domain DB drivers & DSNs
	UserDBDriver  string `json:"user_db_driver"`
	UserDBDSN     string `json:"user_db_dsn"`
	SalesDBDriver string `json:"sales_db_driver"`
	SalesDBDSN    string `json:"sales_db_dsn"`
	// …repeat for each domain…

	AppJWTSecret string `json:"app_jwt_secret"`
	APICertFile  string `json:"api_cert_file"`
	APIKeyFile   string `json:"api_key_file"`
}

// Load loads config from $XDG_CONFIG_HOME/realtorinstall/config.json if present,
// then overrides from ENV vars if set.
func Load() (*Config, error) {
	// 1) JSON file
	cfg := &Config{}
	dir, _ := os.UserConfigDir()
	path := filepath.Join(dir, "realtorinstall", "config.json")
	if b, err := os.ReadFile(path); err == nil {
		json.Unmarshal(b, cfg)
	}
	// 2) Override from ENV
	if v := os.Getenv("USER_DB_DRIVER"); v != "" {
		cfg.UserDBDriver = v
	}
	if v := os.Getenv("USER_DB_DSN"); v != "" {
		cfg.UserDBDSN = v
	}
	if v := os.Getenv("SALES_DB_DRIVER"); v != "" {
		cfg.SalesDBDriver = v
	}
	if v := os.Getenv("SALES_DB_DSN"); v != "" {
		cfg.SalesDBDSN = v
	}
	if v := os.Getenv("APP_JWT_SECRET"); v != "" {
		cfg.AppJWTSecret = v
	}
	// …and so on for the rest…
	return cfg, nil
}
