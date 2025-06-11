// api/repos/sqlite_buyer_repo.go
package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type postgresBuyerRepo struct {
	db *sql.DB
}

func NewSQLiteBuyerRepo(dbPath string) (BuyerRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	schema := `
	CREATE TABLE IF NOT EXISTS buyers (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  tenant_id TEXT NOT NULL,
	  first_name TEXT NOT NULL,
	  last_name TEXT NOT NULL,
	  email TEXT NOT NULL,
	  phone TEXT,
	  created_by TEXT NOT NULL,
	  created_at DATETIME NOT NULL,
	  modified_by TEXT NOT NULL,
	  last_modified DATETIME NOT NULL,
	  deleted INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_buyers_tenant ON buyers(tenant_id);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &sqliteBuyerRepo{db: db}, nil
}

func (r *sqliteBuyerRepo) Create(ctx context.Context, b *models.Buyer) (int64, error) {
	if b.TenantID == "" || b.FirstName == "" || b.LastName == "" || b.Email == "" || b.CreatedBy == "" || b.ModifiedBy == "" {
		return 0, errors.New("missing required fields or tenant/audit info")
	}
	now := time.Now().UTC()
	b.CreatedAt = now
	b.LastModified = now
	query := `
	INSERT INTO buyers (
	  tenant_id, first_name, last_name, email, phone, created_by, created_at, modified_by, last_modified, deleted
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 0);
	`
	res, err := r.db.ExecContext(ctx, query,
		b.TenantID,
		b.FirstName,
		b.LastName,
		b.Email,
		b.Phone,
		b.CreatedBy,
		b.CreatedAt,
		b.ModifiedBy,
		b.LastModified,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *sqliteBuyerRepo) GetByID(ctx context.Context, tenantID string, id int64) (*models.Buyer, error) {
	query := `
	SELECT id, tenant_id, first_name, last_name, email, phone, created_by, created_at, modified_by, last_modified, deleted
	FROM buyers
	WHERE tenant_id = ? AND id = ? AND deleted = 0;
	`
	row := r.db.QueryRowContext(ctx, query, tenantID, id)
	var b models.Buyer
	var deletedInt int
	err := row.Scan(
		&b.ID,
		&b.TenantID,
		&b.FirstName,
		&b.LastName,
		&b.Email,
		&b.Phone,
		&b.CreatedBy,
		&b.CreatedAt,
		&b.ModifiedBy,
		&b.LastModified,
		&deletedInt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	b.Deleted = deletedInt != 0
	return &b, nil
}

func (r *sqliteBuyerRepo) ListAll(ctx context.Context, tenantID string) ([]*models.Buyer, error) {
	query := `
	SELECT id, tenant_id, first_name, last_name, email, phone, created_by, created_at, modified_by, last_modified, deleted
	FROM buyers
	WHERE tenant_id = ? AND deleted = 0;
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
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
			&b.TenantID,
			&b.FirstName,
			&b.LastName,
			&b.Email,
			&b.Phone,
			&b.CreatedBy,
			&b.CreatedAt,
			&b.ModifiedBy,
			&b.LastModified,
			&deletedInt,
		); err != nil {
			return nil, err
		}
		b.Deleted = deletedInt != 0
		out = append(out, &b)
	}
	return out, nil
}

func (r *sqliteBuyerRepo) Update(ctx context.Context, b *models.Buyer) error {
	existing, err := r.GetByID(ctx, b.TenantID, b.ID)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	now := time.Now().UTC()
	b.LastModified = now
	query := `
	UPDATE buyers
	SET first_name = ?, last_name = ?, email = ?, phone = ?, modified_by = ?, last_modified = ?, deleted = ?
	WHERE tenant_id = ? AND id = ?;
	`
	_, err = r.db.ExecContext(ctx, query,
		b.FirstName,
		b.LastName,
		b.Email,
		b.Phone,
		b.ModifiedBy,
		b.LastModified,
		boolToInt(b.Deleted),
		b.TenantID,
		b.ID,
	)
	return err
}

func (r *sqliteBuyerRepo) Delete(ctx context.Context, tenantID string, id int64) error {
	existing, err := r.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	query := `
	UPDATE buyers
	SET deleted = 1, modified_by = ?, last_modified = ?
	WHERE tenant_id = ? AND id = ?;
	`
	_, err = r.db.ExecContext(ctx, query,
		existing.ModifiedBy,
		time.Now().UTC(),
		tenantID,
		id,
	)
	return err
}
