package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type postgresPropertyRepo struct {
	db *sql.DB
}

func NewPostgresPropertyRepo(db *sql.DB) PropertyRepo {
	return &postgresPropertyRepo{db: db}
}

func (r *postgresPropertyRepo) Create(ctx context.Context, p *models.Property) (int64, error) {
	if p.TenantID == "" || p.Address == "" || p.City == "" || p.ZIP == "" || p.CreatedBy == "" || p.ModifiedBy == "" {
		return 0, errors.New("missing required fields or tenant/audit info")
	}
	now := time.Now().UTC()
	p.CreatedAt = now
	p.LastModified = now
	query := `
	INSERT INTO properties (
	  tenant_id, address, city, zip, listing_date, created_by, created_at, modified_by, last_modified, deleted
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 0);
	`
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

func (r *postgresPropertyRepo) GetByID(ctx context.Context, tenantID string, id int64) (*models.Property, error) {
	query := `
	SELECT id, tenant_id, address, city, zip, listing_date, created_by, created_at, modified_by, last_modified, deleted
	FROM properties
	WHERE tenant_id = ? AND id = ? AND deleted = 0;
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
	p.Deleted = deletedInt != 0
	return &p, nil
}

func (r *postgresPropertyRepo) ListAll(ctx context.Context, tenantID string) ([]*models.Property, error) {
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
		p.Deleted = deletedInt != 0
		out = append(out, &p)
	}
	return out, nil
}

func (r *postgresPropertyRepo) Update(ctx context.Context, p *models.Property) error {
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

func (r *postgresPropertyRepo) Delete(ctx context.Context, tenantID string, id int64) error {
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
		existing.ModifiedBy,
		time.Now().UTC(),
		tenantID,
		id,
	)
	return err
}

// SummarizeTopProperties sums all payments by joining installments → payments → properties.
func (r *postgresPropertyRepo) SummarizeTopProperties(ctx context.Context, tenantID string) ([]models.PropertyPaymentVolume, error) {
	query := `
        SELECT p.id            AS property_id,
               SUM(pay.amount) AS total_paid
          FROM payments AS pay
          JOIN installments AS inst ON inst.id = pay.installment_id
          JOIN properties AS p    ON p.id   = inst.property_id
         WHERE pay.tenant_id = ? 
           AND pay.deleted   = 0
           AND inst.deleted  = 0
           AND p.deleted     = 0
         GROUP BY p.id
         ORDER BY total_paid DESC;
    `
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.PropertyPaymentVolume
	for rows.Next() {
		var ppv models.PropertyPaymentVolume
		if err := rows.Scan(&ppv.PropertyID, &ppv.TotalPaidAmount); err != nil {
			return nil, err
		}
		out = append(out, ppv)
	}
	return out, nil
}
