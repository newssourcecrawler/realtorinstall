package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type postgresLocationPricingRepo struct {
	db *sql.DB
}

func NewPostgresLocationPricingRepo(db *sql.DB) LocationPricingRepo {
	return &postgresLocationPricingRepo{db: db}
}

func (r *postgresLocationPricingRepo) Create(ctx context.Context, lp *models.LocationPricing) (int64, error) {
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

func (r *postgresLocationPricingRepo) GetByID(ctx context.Context, tenantID string, id int64) (*models.LocationPricing, error) {
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

func (r *postgresLocationPricingRepo) ListAll(ctx context.Context, tenantID string) ([]*models.LocationPricing, error) {
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

func (r *postgresLocationPricingRepo) Update(ctx context.Context, lp *models.LocationPricing) error {
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

func (r *postgresLocationPricingRepo) Delete(ctx context.Context, tenantID string, id int64) error {
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
