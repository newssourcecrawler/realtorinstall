// internal/repos/sqlite_property_repo.go
package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver

	"github.com/newssourcecrawler/realtorinstall/internal/models"
)

type sqlitePropertyRepo struct {
	db *sql.DB
}

// NewSQLitePropertyRepo opens or creates the SQLite file and ensures schema
func NewSQLitePropertyRepo(dbPath string) (PropertyRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	// Create table if not exists
	schema := `
	CREATE TABLE IF NOT EXISTS properties (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  address TEXT NOT NULL,
	  city TEXT NOT NULL,
	  zip TEXT NOT NULL,
	  location_code TEXT NOT NULL,
	  size_sqft REAL NOT NULL,
	  base_price_usd REAL NOT NULL,
	  listing_date DATETIME NOT NULL,
	  last_modified DATETIME NOT NULL
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &sqlitePropertyRepo{db: db}, nil
}

func (r *sqlitePropertyRepo) Create(ctx context.Context, p *models.Property) (int64, error) {
	now := time.Now().UTC()
	stmt := `INSERT INTO properties
	  (address, city, zip, location_code, size_sqft, base_price_usd, listing_date, last_modified)
	  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, stmt,
		p.Address, p.City, p.ZIP, p.LocationCode,
		p.SizeSqFt, p.BasePriceUSD, p.ListingDate, now,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *sqlitePropertyRepo) GetByID(ctx context.Context, id int64) (*models.Property, error) {
	row := r.db.QueryRowContext(ctx, `
	  SELECT id, address, city, zip, location_code, size_sqft, base_price_usd, listing_date, last_modified
	  FROM properties WHERE id = ?`, id)

	var p models.Property
	var listingDate, lastMod string
	if err := row.Scan(
		&p.ID, &p.Address, &p.City, &p.ZIP, &p.LocationCode,
		&p.SizeSqFt, &p.BasePriceUSD, &listingDate, &lastMod,
	); errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	p.ListingDate, _ = time.Parse(time.RFC3339, listingDate)
	p.LastModified, _ = time.Parse(time.RFC3339, lastMod)
	return &p, nil
}

func (r *sqlitePropertyRepo) ListAll(ctx context.Context) ([]*models.Property, error) {
	rows, err := r.db.QueryContext(ctx, `
	  SELECT id, address, city, zip, location_code, size_sqft, base_price_usd, listing_date, last_modified
	  FROM properties ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var props []*models.Property
	for rows.Next() {
		var p models.Property
		var listingDate, lastMod string
		if err := rows.Scan(
			&p.ID, &p.Address, &p.City, &p.ZIP, &p.LocationCode,
			&p.SizeSqFt, &p.BasePriceUSD, &listingDate, &lastMod,
		); err != nil {
			return nil, err
		}
		p.ListingDate, _ = time.Parse(time.RFC3339, listingDate)
		p.LastModified, _ = time.Parse(time.RFC3339, lastMod)
		props = append(props, &p)
	}
	return props, nil
}

func (r *sqlitePropertyRepo) Update(ctx context.Context, p *models.Property) error {
	p.LastModified = time.Now().UTC()
	stmt := `
	  UPDATE properties SET
	    address = ?, city = ?, zip = ?, location_code = ?, size_sqft = ?, base_price_usd = ?, listing_date = ?, last_modified = ?
	  WHERE id = ?`
	_, err := r.db.ExecContext(ctx, stmt,
		p.Address, p.City, p.ZIP, p.LocationCode,
		p.SizeSqFt, p.BasePriceUSD, p.ListingDate, p.LastModified,
		p.ID,
	)
	return err
}

func (r *sqlitePropertyRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM properties WHERE id = ?`, id)
	return err
}
