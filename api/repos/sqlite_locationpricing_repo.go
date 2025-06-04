// api/repos/sqlite_locationpricing_repo.go
package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type sqliteLocationPricingRepo struct {
	db *sql.DB
}

func NewSQLiteLocationPricingRepo(dbPath string) (LocationPricingRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	schema := `
	CREATE TABLE IF NOT EXISTS location_pricing (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  tenant_id TEXT NOT NULL,
	  zip_code TEXT NOT NULL,
	  city TEXT NOT NULL,
	  price_per_sqft REAL NOT NULL,
	  effective_date DATETIME NOT NULL,
	  created_by TEXT NOT NULL,
	  created_at DATETIME NOT NULL,
	  modified_by TEXT NOT NULL,
	  last_modified DATETIME NOT NULL,
	  deleted INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_locationpricing_tenant ON location_pricing(tenant_id);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &sqliteLocationPricingRepo{db: db}, nil
}

func (r *sqliteLocationPricingRepo) Create(ctx context.Context, lp *models.LocationPricing) (int64, error) {
	if lp.TenantID == "" || lp.ZipCode == "" || lp.City == "" || lp.CreatedBy == "" || lp.ModifiedBy == "" || lp.PricePerSqFt <= 0 || lp.EffectiveDate.IsZero() {
		return 0, errors.New("missing required fields or tenant/audit info")
	}
	now := time.Now().UTC()
	lp.CreatedAt = now
	lp.LastModified = now
	query := `
	INSERT INTO location_pricing (
	  tenant_id, zip_code, city, price_per_sqft, effective_date,
	  created_by, created_at, modified_by, last_modified, deleted
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 0);
	`
	res, err := r.db.ExecContext(ctx, query,
		lp.TenantID,
		lp.ZipCode,
		lp.City,
		lp.PricePerSqFt,
		lp.EffectiveDate,
		lp.CreatedBy,
		lp.CreatedAt,
		lp.ModifiedBy,
		lp.LastModified,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *sqliteLocationPricingRepo) GetByID(ctx context.Context, tenantID string, id int64) (*models.LocationPricing, error) {
	query := `
	SELECT id, tenant_id, zip_code, city, price_per_sqft, effective_date, created_by, created_at, modified_by, last_modified, deleted
	FROM location_pricing
	WHERE tenant_id = ? AND id = ? AND deleted = 0;
	`
	row := r.db.QueryRowContext(ctx, query, tenantID, id)
	var lp models.LocationPricing
	var deletedInt int
	err := row.Scan(
		&lp.ID,
		&lp.TenantID,
		&lp.ZipCode,
		&lp.City,
		&lp.PricePerSqFt,
		&lp.EffectiveDate,
		&lp.CreatedBy,
		&lp.CreatedAt,
		&lp.ModifiedBy,
		&lp.LastModified,
		&deletedInt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	lp.Deleted = deletedInt != 0
	return &lp, nil
}

func (r *sqliteLocationPricingRepo) ListAll(ctx context.Context, tenantID string) ([]*models.LocationPricing, error) {
	query := `
	SELECT id, tenant_id, zip_code, city, price_per_sqft, effective_date, created_by, created_at, modified_by, last_modified, deleted
	FROM location_pricing
	WHERE tenant_id = ? AND deleted = 0;
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
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
			&lp.TenantID,
			&lp.ZipCode,
			&lp.City,
			&lp.PricePerSqFt,
			&lp.EffectiveDate,
			&lp.CreatedBy,
			&lp.CreatedAt,
			&lp.ModifiedBy,
			&lp.LastModified,
			&deletedInt,
		); err != nil {
			return nil, err
		}
		lp.Deleted = deletedInt != 0
		out = append(out, &lp)
	}
	return out, nil
}

func (r *sqliteLocationPricingRepo) Update(ctx context.Context, lp *models.LocationPricing) error {
	existing, err := r.GetByID(ctx, lp.TenantID, lp.ID)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	now := time.Now().UTC()
	lp.LastModified = now
	query := `
	UPDATE location_pricing
	SET zip_code = ?, city = ?, price_per_sqft = ?, effective_date = ?, modified_by = ?, last_modified = ?, deleted = ?
	WHERE tenant_id = ? AND id = ?;
	`
	_, err = r.db.ExecContext(ctx, query,
		lp.ZipCode,
		lp.City,
		lp.PricePerSqFt,
		lp.EffectiveDate,
		lp.ModifiedBy,
		lp.LastModified,
		boolToInt(lp.Deleted),
		lp.TenantID,
		lp.ID,
	)
	return err
}

func (r *sqliteLocationPricingRepo) Delete(ctx context.Context, tenantID string, id int64) error {
	existing, err := r.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	query := `
	UPDATE location_pricing
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
