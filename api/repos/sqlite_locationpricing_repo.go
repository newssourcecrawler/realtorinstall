package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
)

// sqliteLocationPricingRepo implements LocationPricingRepo using SQLite.
type sqliteLocationPricingRepo struct {
	db *sql.DB
}

// NewSQLiteLocationPricingRepo opens/creates "location_pricing" table.
func NewSQLiteLocationPricingRepo(dbPath string) (LocationPricingRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// We add a 'deleted' column for soft‐deletes
	schema := `
	CREATE TABLE IF NOT EXISTS location_pricing (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  zip_code TEXT NOT NULL,
	  city TEXT NOT NULL,
	  price_per_sqft REAL NOT NULL,
	  effective_date DATETIME NOT NULL,
	  created_at DATETIME NOT NULL,
	  last_modified DATETIME NOT NULL,
	  deleted INTEGER NOT NULL DEFAULT 0
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &sqliteLocationPricingRepo{db: db}, nil
}

// Create inserts a new LocationPricing record.
func (r *sqliteLocationPricingRepo) Create(ctx context.Context, lp *models.LocationPricing) (int64, error) {
	query := `
	INSERT INTO location_pricing (
	  zip_code,
	  city,
	  price_per_sqft,
	  effective_date,
	  created_at,
	  last_modified,
	  deleted
	) VALUES (?, ?, ?, ?, ?, ?, 0);
	`
	// Ensure CreatedAt/LastModified are set
	if lp.EffectiveDate.IsZero() {
		return 0, errors.New("effective_date is required")
	}
	now := time.Now().UTC()
	lp.CreatedAt = now
	lp.LastModified = now

	res, err := r.db.ExecContext(ctx, query,
		lp.ZipCode,
		lp.City,
		lp.PricePerSqFt,
		lp.EffectiveDate,
		lp.CreatedAt,
		lp.LastModified,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetByID retrieves one LocationPricing, even if soft‐deleted.
func (r *sqliteLocationPricingRepo) GetByID(ctx context.Context, id int64) (*models.LocationPricing, error) {
	query := `
	SELECT id, zip_code, city, price_per_sqft, effective_date, created_at, last_modified, deleted
	FROM location_pricing
	WHERE id = ?;
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var lp models.LocationPricing
	var deletedInt int
	err := row.Scan(
		&lp.ID,
		&lp.ZipCode,
		&lp.City,
		&lp.PricePerSqFt,
		&lp.EffectiveDate,
		&lp.CreatedAt,
		&lp.LastModified,
		&deletedInt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	// We store deleted but services can decide whether to ignore it
	return &lp, nil
}

// ListAll returns only non‐deleted LocationPricing rows.
func (r *sqliteLocationPricingRepo) ListAll(ctx context.Context) ([]*models.LocationPricing, error) {
	query := `
	SELECT id, zip_code, city, price_per_sqft, effective_date, created_at, last_modified, deleted
	FROM location_pricing
	WHERE deleted = 0;
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.LocationPricing
	for rows.Next() {
		var lp models.LocationPricing
		var deletedInt int
		if err := rows.Scan(
			&lp.ID,
			&lp.ZipCode,
			&lp.City,
			&lp.PricePerSqFt,
			&lp.EffectiveDate,
			&lp.CreatedAt,
			&lp.LastModified,
			&deletedInt,
		); err != nil {
			return nil, err
		}
		out = append(out, &lp)
	}
	return out, nil
}

// Update modifies an existing LocationPricing (including its deleted flag).
func (r *sqliteLocationPricingRepo) Update(ctx context.Context, lp *models.LocationPricing) error {
	query := `
	UPDATE location_pricing
	SET
	  zip_code = ?,
	  city = ?,
	  price_per_sqft = ?,
	  effective_date = ?,
	  last_modified = ?,
	  deleted = ?
	WHERE id = ?;
	`
	now := time.Now().UTC()
	lp.LastModified = now

	_, err := r.db.ExecContext(ctx, query,
		lp.ZipCode,
		lp.City,
		lp.PricePerSqFt,
		lp.EffectiveDate,
		lp.LastModified,
		boolToInt(lp.Deleted),
		lp.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

// Delete performs a soft‐delete by setting deleted = 1.
func (r *sqliteLocationPricingRepo) Delete(ctx context.Context, id int64) error {
	// First, check that the row exists
	_, err := r.GetByID(ctx, id)
	if err == repos.ErrNotFound {
		return repos.ErrNotFound
	} else if err != nil {
		return err
	}

	query := `
	UPDATE location_pricing
	SET deleted = 1, last_modified = ?
	WHERE id = ?;
	`
	_, err = r.db.ExecContext(ctx, query, time.Now().UTC(), id)
	return err
}
