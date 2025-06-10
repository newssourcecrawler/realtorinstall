// api/internal/config/config.go
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	// Per-domain DB drivers & DSNs
	UserDBDriver           string `json:"user_db_driver"`
	UserDBDSN              string `json:"user_db_dsn"`
	SalesDBDriver          string `json:"sales_db_driver"`
	SalesDBDSN             string `json:"sales_db_dsn"`
	BuyerDBDriver          string `json:"buyer_db_driver"`
	BuyerDBDSN             string `json:"buyer_db_dsn"`
	CommissionDBDriver     string `json:"commission_db_driver"`
	CommissionDBDSN        string `json:"commission_db_dsn"`
	InstDBDriver           string `json:"Installment_db_driver"`
	InstDBDSN              string `json:"Installment_db_dsn"`
	PlanDBDriver           string `json:"plan_db_driver"`
	PlanDBDSN              string `json:"plan_db_dsn"`
	IntroDBDriver          string `json:"introductions_db_driver"`
	IntroDBDSN             string `json:"introductions_db_dsn"`
	LettingsDBDriver       string `json:"lettings_db_driver"`
	LettingsDBDSN          string `json:"lettings_db_dsn"`
	PricingDBDriver        string `json:"pricing_db_driver"`
	PricingDBDSN           string `json:"pricing_db_dsn"`
	PayDBDriver            string `json:"pay_db_driver"`
	PayDBDSN               string `json:"pay_db_dsn"`
	PermissionDBDriver     string `json:"permissions_db_driver"`
	PermissionDBDSN        string `json:"permissions_db_dsn"`
	PropertyDBDriver       string `json:"property_db_driver"`
	PropertyDBDSN          string `json:"property_db_dsn"`
	RolePermissionDBDriver string `json:"rolepermissions_db_driver"`
	RolePermissionDBDSN    string `json:"rolepermissions_db_dsn"`
	RoleDBDriver           string `json:"role_db_driver"`
	RoleDBDSN              string `json:"role_db_dsn"`
	UserRoleDBDriver       string `json:"userrole_db_driver"`
	UserRoleDBDSN          string `json:"userrole_db_dsn"`

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
	if v := os.Getenv("COMMISSION_DB_DRIVER"); v != "" {
		cfg.CommissionDBDriver = v
	}
	if v := os.Getenv("COMMISSION_DB_DSN"); v != "" {
		cfg.CommissionDBDSN = v
	}
	if v := os.Getenv("PROPERTY_DB_DRIVER"); v != "" {
		cfg.PropertyDBDriver = v
	}
	if v := os.Getenv("PROPERTY_DB_DSN"); v != "" {
		cfg.PropertyDBDSN = v
	}
	if v := os.Getenv("PRICING_DB_DRIVER"); v != "" {
		cfg.PricingDBDriver = v
	}
	if v := os.Getenv("PRICING_DB_DSN"); v != "" {
		cfg.PricingDBDSN = v
	}
	if v := os.Getenv("BUYER_DB_DRIVER"); v != "" {
		cfg.BuyerDBDriver = v
	}
	if v := os.Getenv("BUYER_DB_DSN"); v != "" {
		cfg.BuyerDBDSN = v
	}
	if v := os.Getenv("PLAN_DB_DRIVER"); v != "" {
		cfg.PlanDBDriver = v
	}
	if v := os.Getenv("PLAN_DB_DSN"); v != "" {
		cfg.PlanDBDSN = v
	}
	if v := os.Getenv("INSTALLMENTS_DB_DRIVER"); v != "" {
		cfg.InstDBDriver = v
	}
	if v := os.Getenv("INSTALLMENTS_DB_DSN"); v != "" {
		cfg.InstDBDSN = v
	}
	if v := os.Getenv("PAYMENTS_DB_DRIVER"); v != "" {
		cfg.PayDBDriver = v
	}
	if v := os.Getenv("PAYMENTS_DB_DSN"); v != "" {
		cfg.PayDBDSN = v
	}
	if v := os.Getenv("INTRODUCTIONS_DB_DRIVER"); v != "" {
		cfg.IntroDBDriver = v
	}
	if v := os.Getenv("INTRODUCTIONS_DB_DSN"); v != "" {
		cfg.IntroDBDSN = v
	}
	if v := os.Getenv("LETTINGS_DB_DRIVER"); v != "" {
		cfg.LettingsDBDriver = v
	}
	if v := os.Getenv("LETTINGS_DB_DSN"); v != "" {
		cfg.LettingsDBDSN = v
	}
	if v := os.Getenv("ROLEPERMISSION_DB_DRIVER"); v != "" {
		cfg.RolePermissionDBDriver = v
	}
	if v := os.Getenv("ROLEPERMISSION_DB_DSN"); v != "" {
		cfg.RolePermissionDBDSN = v
	}
	if v := os.Getenv("ROLE_DB_DRIVER"); v != "" {
		cfg.RoleDBDriver = v
	}
	if v := os.Getenv("ROLE_DB_DSN"); v != "" {
		cfg.RoleDBDSN = v
	}
	if v := os.Getenv("PERMISSIONS_DB_DRIVER"); v != "" {
		cfg.PermissionDBDriver = v
	}
	if v := os.Getenv("PERMISSIONS_DB_DSN"); v != "" {
		cfg.PermissionDBDSN = v
	}
	if v := os.Getenv("USERROLE_DB_DRIVER"); v != "" {
		cfg.UserRoleDBDriver = v
	}
	if v := os.Getenv("USERROLE_DB_DSN"); v != "" {
		cfg.UserRoleDBDSN = v
	}

	if v := os.Getenv("APP_JWT_SECRET"); v != "" {
		cfg.AppJWTSecret = v
	}
	if v := os.Getenv("APP_CERT_FILE"); v != "" {
		cfg.APICertFile = v
	}
	if v := os.Getenv("APP_KEY_FILE"); v != "" {
		cfg.APIKeyFile = v
	}
	return cfg, nil
}
