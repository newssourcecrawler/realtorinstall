package repos

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/newssourcecrawler/realtorinstall/api/models"
	"github.com/newssourcecrawler/realtorinstall/api/repos"
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

	// Create a "properties" table (with soft-delete support).
	schema := `
	CREATE TABLE IF NOT EXISTS properties (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  address TEXT NOT NULL,
	  city TEXT NOT NULL,
	  zip TEXT NOT NULL,
	  listing_date DATETIME NOT NULL,
	  created_at DATETIME NOT NULL,
	  last_modified DATETIME NOT NULL,
	  deleted INTEGER NOT NULL DEFAULT 0
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &sqlitePropertyRepo{db: db}, nil
}

// Create inserts a new Property.
func (r *sqlitePropertyRepo) Create(ctx context.Context, p *models.Property) (int64, error) {
	if p.Address == "" || p.City == "" || p.ZIP == "" {
		//return 0, errors.New("address, city, and ZIP are required")
		return 0, repos.AddrNotFound
	}

	query := `
	INSERT INTO properties (address, city, zip, listing_date, created_at, last_modified, deleted)
	VALUES (?, ?, ?, ?, ?, ?, ?);
	`
	res, err := r.db.ExecContext(ctx, query,
		p.Address,
		p.City,
		p.ZIP,
		p.ListingDate,
		p.CreatedAt,
		p.LastModified,
		boolToInt(p.Deleted),
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetByID retrieves a Property by ID (even if soft-deleted).
func (r *sqlitePropertyRepo) GetByID(ctx context.Context, id int64) (*models.Property, error) {
	query := `
	SELECT id, address, city, zip, listing_date, created_at, last_modified, deleted
	FROM properties
	WHERE id = ?;
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var p models.Property
	var deletedInt int
	err := row.Scan(
		&p.ID,
		&p.Address,
		&p.City,
		&p.ZIP,
		&p.ListingDate,
		&p.CreatedAt,
		&p.LastModified,
		&deletedInt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	p.Deleted = intToBool(deletedInt)
	return &p, nil
}

// ListAll returns all non-deleted Properties.
func (r *sqlitePropertyRepo) ListAll(ctx context.Context) ([]*models.Property, error) {
	query := `
	SELECT id, address, city, zip, listing_date, created_at, last_modified, deleted
	FROM properties
	WHERE deleted = 0;
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Property
	for rows.Next() {
		var p models.Property
		var deletedInt int
		if err := rows.Scan(
			&p.ID,
			&p.Address,
			&p.City,
			&p.ZIP,
			&p.ListingDate,
			&p.CreatedAt,
			&p.LastModified,
			&deletedInt,
		); err != nil {
			return nil, err
		}
		p.Deleted = intToBool(deletedInt)
		out = append(out, &p)
	}
	return out, nil
}

// Update modifies an existing Property (including its Deleted flag).
func (r *sqlitePropertyRepo) Update(ctx context.Context, p *models.Property) error {
	// Check existence first
	_, err := r.GetByID(ctx, p.ID)
	if err != nil {
		return err
	}

	query := `
	UPDATE properties
	SET address = ?, city = ?, zip = ?, listing_date = ?, last_modified = ?, deleted = ?
	WHERE id = ?;
	`
	_, err = r.db.ExecContext(ctx, query,
		p.Address,
		p.City,
		p.ZIP,
		p.ListingDate,
		p.LastModified,
		boolToInt(p.Deleted),
		p.ID,
	)
	return err
}

// Delete performs a soft-delete by setting deleted = 1.
func (r *sqlitePropertyRepo) Delete(ctx context.Context, id int64) error {
	// Confirm existence
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	query := `
	UPDATE properties
	SET deleted = 1, last_modified = ?
	WHERE id = ?;
	`
	_, err = r.db.ExecContext(ctx, query, nowUTCString(), id)
	return err
}
