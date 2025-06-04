package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type sqlitePropertyRepo struct {
	db *sql.DB
}

// NewSQLitePropertyRepo opens/creates the properties table (with tenant_id, etc.).
func NewSQLitePropertyRepo(dbPath string) (PropertyRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	schema := `
    CREATE TABLE IF NOT EXISTS properties (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      tenant_id TEXT NOT NULL,
      address TEXT NOT NULL,
      city TEXT NOT NULL,
      zip TEXT NOT NULL,
      listing_date DATETIME NOT NULL,
      created_by TEXT NOT NULL,
      created_at DATETIME NOT NULL,
      modified_by TEXT NOT NULL,
      last_modified DATETIME NOT NULL,
      deleted INTEGER NOT NULL DEFAULT 0
    );
    `
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	// Ensure index on tenant_id to speed up per‐tenant queries
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_properties_tenant ON properties(tenant_id);`); err != nil {
		return nil, err
	}
	return &sqlitePropertyRepo{db: db}, nil
}

func (r *sqlitePropertyRepo) Create(ctx context.Context, p *models.Property) (int64, error) {
	if p.TenantID == "" || p.Address == "" || p.City == "" || p.ZIP == "" || p.CreatedBy == "" || p.ModifiedBy == "" {
		return 0, errors.New("missing required fields or tenant_id/created_by/modified_by")
	}
	query := `
    INSERT INTO properties (
      tenant_id, address, city, zip, listing_date, created_by, created_at, modified_by, last_modified, deleted
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 0);
    `
	now := time.Now().UTC()
	p.CreatedAt = now
	p.LastModified = now

	res, err := r.db.ExecContext(ctx, query,
		p.TenantID,
		p.Address,
		p.City,
		p.ZIP,
		p.ListingDate,
		p.CreatedBy,
		p.CreatedAt,
		p.ModifiedBy,
		p.LastModified,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *sqlitePropertyRepo) GetByID(ctx context.Context, tenantID string, id int64) (*models.Property, error) {
	query := `
    SELECT id, tenant_id, address, city, zip, listing_date, created_by, created_at, modified_by, last_modified, deleted
    FROM properties
    WHERE tenant_id = ? AND id = ?;
    `
	row := r.db.QueryRowContext(ctx, query, tenantID, id)

	var p models.Property
	var deletedInt int
	err := row.Scan(
		&p.ID,
		&p.TenantID,
		&p.Address,
		&p.City,
		&p.ZIP,
		&p.ListingDate,
		&p.CreatedBy,
		&p.CreatedAt,
		&p.ModifiedBy,
		&p.LastModified,
		&deletedInt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	p.Deleted = (deletedInt != 0)
	return &p, nil
}

func (r *sqlitePropertyRepo) ListAll(ctx context.Context, tenantID string) ([]*models.Property, error) {
	query := `
    SELECT id, tenant_id, address, city, zip, listing_date, created_by, created_at, modified_by, last_modified, deleted
    FROM properties
    WHERE tenant_id = ? AND deleted = 0;
    `
	rows, err := r.db.QueryContext(ctx, query, tenantID)
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
			&p.TenantID,
			&p.Address,
			&p.City,
			&p.ZIP,
			&p.ListingDate,
			&p.CreatedBy,
			&p.CreatedAt,
			&p.ModifiedBy,
			&p.LastModified,
			&deletedInt,
		); err != nil {
			return nil, err
		}
		p.Deleted = (deletedInt != 0)
		out = append(out, &p)
	}
	return out, nil
}

func (r *sqlitePropertyRepo) Update(ctx context.Context, p *models.Property) error {
	// Must confirm that p.TenantID and p.ID match an existing row
	existing, err := r.GetByID(ctx, p.TenantID, p.ID)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}

	now := time.Now().UTC()
	p.LastModified = now

	query := `
    UPDATE properties
    SET address = ?, city = ?, zip = ?, listing_date = ?, modified_by = ?, last_modified = ?, deleted = ?
    WHERE tenant_id = ? AND id = ?;
    `
	_, err = r.db.ExecContext(ctx, query,
		p.Address,
		p.City,
		p.ZIP,
		p.ListingDate,
		p.ModifiedBy,
		p.LastModified,
		boolToInt(p.Deleted),
		p.TenantID,
		p.ID,
	)
	return err
}

func (r *sqlitePropertyRepo) Delete(ctx context.Context, tenantID string, id int64) error {
	// Soft‐delete by setting deleted = 1
	existing, err := r.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	query := `
    UPDATE properties
    SET deleted = 1, modified_by = ?, last_modified = ?
    WHERE tenant_id = ? AND id = ?;
    `
	_, err = r.db.ExecContext(ctx, query,
		existing.ModifiedBy, // or pass a new “deleter” if you track delete‐by separately
		time.Now().UTC(),
		tenantID,
		id,
	)
	return err
}
