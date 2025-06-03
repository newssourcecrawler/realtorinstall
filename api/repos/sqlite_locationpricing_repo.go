package repos

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/newssourcecrawler/realtorinstall/internal/models"
)

type sqliteLocationPricingRepo struct{ db *sql.DB }

// NewSQLiteLocationPricingRepo opens/creates "location_pricing" table
func NewSQLiteLocationPricingRepo(dbPath string) (LocationPricingRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	schema := `
	CREATE TABLE IF NOT EXISTS location_pricing (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  zip_code TEXT NOT NULL,
	  city TEXT NOT NULL,
	  price_per_sqft REAL NOT NULL,
	  effective_date DATETIME NOT NULL,
	  created_at DATETIME NOT NULL,
	  last_modified DATETIME NOT NULL
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &sqliteLocationPricingRepo{db: db}, nil
}

func (r *sqliteLocationPricingRepo) Create(ctx context.Context, lp *models.LocationPricing) (int64, error) {
	return 0, nil
}
func (r *sqliteLocationPricingRepo) GetByID(ctx context.Context, id int64) (*models.LocationPricing, error) {
	return nil, nil
}
func (r *sqliteLocationPricingRepo) ListAll(ctx context.Context) ([]*models.LocationPricing, error) {
	return nil, nil
}
func (r *sqliteLocationPricingRepo) Update(ctx context.Context, lp *models.LocationPricing) error {
	return nil
}
func (r *sqliteLocationPricingRepo) Delete(ctx context.Context, id int64) error {
	return nil
}
