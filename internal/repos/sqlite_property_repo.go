package repos

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/newssourcecrawler/realtorinstall/internal/models"
)

// sqlitePropertyRepo implements PropertyRepo using SQLite.
type sqlitePropertyRepo struct {
	db *sql.DB
}

// NewSQLitePropertyRepo opens (or creates) the SQLite file and ensures the "properties" table exists.
func NewSQLitePropertyRepo(dbPath string) (PropertyRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create a "properties" table that matches models.Property (string dates)
	schema := `
	CREATE TABLE IF NOT EXISTS properties (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  address TEXT NOT NULL,
	  city TEXT NOT NULL,
	  zip TEXT NOT NULL,
	  listing_date TEXT NOT NULL,
	  created_at TEXT NOT NULL,
	  last_modified TEXT NOT NULL
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &sqlitePropertyRepo{db: db}, nil
}

// Create inserts a new Property. It expects p.ListingDate, p.CreatedAt, p.LastModified as RFC3339 strings.
func (r *sqlitePropertyRepo) Create(ctx context.Context, p *models.Property) (int64, error) {
	query := `
	INSERT INTO properties (address, city, zip, listing_date, created_at, last_modified)
	VALUES (?, ?, ?, ?, ?, ?);
	`
	res, err := r.db.ExecContext(ctx, query,
		p.Address,
		p.City,
		p.ZIP,
		p.ListingDate,
		p.CreatedAt,
		p.LastModified,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetByID retrieves a Property by ID.
func (r *sqlitePropertyRepo) GetByID(ctx context.Context, id int64) (*models.Property, error) {
	query := `
	SELECT id, address, city, zip, listing_date, created_at, last_modified
	FROM properties WHERE id = ?;
	`
	row := r.db.QueryRowContext(ctx, query, id)
	var p models.Property
	err := row.Scan(
		&p.ID,
		&p.Address,
		&p.City,
		&p.ZIP,
		&p.ListingDate,
		&p.CreatedAt,
		&p.LastModified,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ListAll returns all Properties.
func (r *sqlitePropertyRepo) ListAll(ctx context.Context) ([]*models.Property, error) {
	query := `
	SELECT id, address, city, zip, listing_date, created_at, last_modified
	FROM properties;
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var props []*models.Property
	for rows.Next() {
		var p models.Property
		if err := rows.Scan(
			&p.ID,
			&p.Address,
			&p.City,
			&p.ZIP,
			&p.ListingDate,
			&p.CreatedAt,
			&p.LastModified,
		); err != nil {
			return nil, err
		}
		props = append(props, &p)
	}
	return props, nil
}

// Update modifies an existing Property.
func (r *sqlitePropertyRepo) Update(ctx context.Context, p *models.Property) error {
	query := `
	UPDATE properties
	SET address = ?, city = ?, zip = ?, listing_date = ?, created_at = ?, last_modified = ?
	WHERE id = ?;
	`
	_, err := r.db.ExecContext(ctx, query,
		p.Address,
		p.City,
		p.ZIP,
		p.ListingDate,
		p.CreatedAt,
		p.LastModified,
		p.ID,
	)
	return err
}

// Delete removes a Property by ID.
func (r *sqlitePropertyRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM properties WHERE id = ?;`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
