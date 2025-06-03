// internal/utils/helpers.go
package utils

import (
	"encoding/json"
	"os"
)

// Config holds application‐wide settings
type Config struct {
	DatabasePath string `json:"DatabasePath"`
	SMTPHost     string `json:"SMTPHost"`
	SMTPPort     int    `json:"SMTPPort"`
	SMTPUser     string `json:"SMTPUser"`
	SMTPPass     string `json:"SMTPPass"`
	AdminUser    string `json:"AdminUser"`
	AdminHash    string `json:"AdminHash"`
}

// LoadConfig reads config from the given JSON file
func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
