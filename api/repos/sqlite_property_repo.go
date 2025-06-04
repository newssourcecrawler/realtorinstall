package repos

import (
	"context"
	"database/sql"
	"strconv"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/newssourcecrawler/realtorinstall/api/models"
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

	// Create a "properties" table that matches models.Property (including 'deleted' flag)
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

// GetByID retrieves a Property by ID (soft‐deleted rows are still returned so service can mark them).
func (r *sqlitePropertyRepo) GetByID(ctx context.Context, id string) (*models.Property, error) {
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, ErrNotFound
	}

	query := `
	SELECT id, address, city, zip, listing_date, created_at, last_modified, deleted
	FROM properties
	WHERE id = ?;
	`
	row := r.db.QueryRowContext(ctx, query, intID)

	var p models.Property
	var deletedInt int
	err = row.Scan(
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

// ListAll returns all non‐deleted Properties.
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

	var props []*models.Property
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
		props = append(props, &p)
	}
	return props, nil
}

// Update modifies an existing Property (including its Deleted flag).
func (r *sqlitePropertyRepo) Update(ctx context.Context, p *models.Property) error {
	query := `
	UPDATE properties
	SET address = ?, city = ?, zip = ?, listing_date = ?, last_modified = ?, deleted = ?
	WHERE id = ?;
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
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

func (r *sqlitePropertyRepo) Delete(ctx context.Context, id string) error {
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return ErrNotFound
	}
	result, err := r.db.ExecContext(ctx,
		`UPDATE properties
         SET deleted = 1, last_modified = CURRENT_TIMESTAMP
         WHERE id = ?;`, intID)
	if err != nil {
		return err // some DB error
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
