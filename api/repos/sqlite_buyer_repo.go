package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/newssourcecrawler/realtorinstall/api/models"
)

// sqliteBuyerRepo implements BuyerRepo using SQLite.
type sqliteBuyerRepo struct {
	db *sql.DB
}

// NewSQLiteBuyerRepo opens (or creates) the SQLite file and ensures the "buyers" table exists.
func NewSQLiteBuyerRepo(dbPath string) (BuyerRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS buyers (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  first_name TEXT NOT NULL,
	  last_name TEXT NOT NULL,
	  email TEXT NOT NULL,
	  phone TEXT,
	  created_at DATETIME NOT NULL,
	  last_modified DATETIME NOT NULL,
	  deleted INTEGER NOT NULL DEFAULT 0
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &sqliteBuyerRepo{db: db}, nil
}

// Create inserts a new Buyer.
func (r *sqliteBuyerRepo) Create(ctx context.Context, b *models.Buyer) (int64, error) {
	if b.FirstName == "" || b.LastName == "" || b.Email == "" {
		return 0, errors.New("first name, last name and email are required")
		//return 0, repos.NameEmailNotFound
	}

	query := `
	INSERT INTO buyers (first_name, last_name, email, phone, created_at, last_modified, deleted)
	VALUES (?, ?, ?, ?, ?, ?, ?);
	`
	res, err := r.db.ExecContext(ctx, query,
		b.FirstName,
		b.LastName,
		b.Email,
		b.Phone,
		b.CreatedAt,
		b.LastModified,
		boolToInt(b.Deleted),
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetByID retrieves a Buyer by ID (even if soft-deleted).
func (r *sqliteBuyerRepo) GetByID(ctx context.Context, id int64) (*models.Buyer, error) {
	query := `
	SELECT id, first_name, last_name, email, phone, created_at, last_modified, deleted
	FROM buyers
	WHERE id = ?;
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var b models.Buyer
	var deletedInt int
	err := row.Scan(
		&b.ID,
		&b.FirstName,
		&b.LastName,
		&b.Email,
		&b.Phone,
		&b.CreatedAt,
		&b.LastModified,
		&deletedInt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	b.Deleted = intToBool(deletedInt)
	return &b, nil
}

// ListAll returns all non-deleted Buyers.
func (r *sqliteBuyerRepo) ListAll(ctx context.Context) ([]*models.Buyer, error) {
	query := `
	SELECT id, first_name, last_name, email, phone, created_at, last_modified, deleted
	FROM buyers
	WHERE deleted = 0;
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Buyer
	for rows.Next() {
		var b models.Buyer
		var deletedInt int
		if err := rows.Scan(
			&b.ID,
			&b.FirstName,
			&b.LastName,
			&b.Email,
			&b.Phone,
			&b.CreatedAt,
			&b.LastModified,
			&deletedInt,
		); err != nil {
			return nil, err
		}
		b.Deleted = intToBool(deletedInt)
		out = append(out, &b)
	}
	return out, nil
}

// Update modifies an existing Buyer (including the Deleted flag).
func (r *sqliteBuyerRepo) Update(ctx context.Context, id int64, b *models.Buyer) error {
	// Check existence first (optional; service may already have done this).
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	query := `
	UPDATE buyers
	SET first_name = ?, last_name = ?, email = ?, phone = ?, last_modified = ?, deleted = ?
	WHERE id = ?;
	`
	_, err = r.db.ExecContext(ctx, query,
		b.FirstName,
		b.LastName,
		b.Email,
		b.Phone,
		b.LastModified,
		boolToInt(b.Deleted),
		id,
	)
	return err
}

// Delete performs a soft-delete (sets deleted = 1).
func (r *sqliteBuyerRepo) Delete(ctx context.Context, id int64) error {
	// Confirm existence
	_, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}
	query := `
	UPDATE buyers
	SET deleted = 1, last_modified = ?
	WHERE id = ?;
	`
	_, err = r.db.ExecContext(ctx, query, nowUTCString(), id)
	return err
}

// Helper: convert 0/1 → bool
func intToBool(i int) bool {
	return i != 0
}

// nowUTCString returns current UTC time formatted for SQLite.
func nowUTCString() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05")
}
