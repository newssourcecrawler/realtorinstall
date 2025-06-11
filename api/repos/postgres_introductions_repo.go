package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type postgresIntroductionsRepo struct {
	db *sql.DB
}

func NewPostgresIntroductionsRepo(db *sql.DB) IntroductionsRepo {
	return &postgresIntroductionsRepo{db: db}
}

func (r *postgresIntroductionsRepo) Create(ctx context.Context, intro *models.Introductions) (int64, error) {
	if intro.TenantID == "" ||
		intro.IntroducerID == 0 ||
		intro.IntroducedParty == "" ||
		intro.PropertyID == 0 ||
		intro.IntroDate.IsZero() ||
		intro.FeeType == "" ||
		intro.CreatedBy == "" ||
		intro.ModifiedBy == "" {
		return 0, errors.New("missing required fields or tenant/audit info")
	}
	now := time.Now().UTC()
	intro.CreatedAt = now
	intro.LastModified = now

	query := `
	INSERT INTO introductions (
	  tenant_id, introducer_id, introduced_party, property_id, transaction_id, intro_date, agreed_fee, fee_type,
	  created_by, created_at, modified_by, last_modified, deleted
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0);
	`
	res, err := r.db.ExecContext(ctx, query,
		intro.TenantID,
		intro.IntroducerID,
		intro.IntroducedParty,
		intro.PropertyID,
		intro.TransactionID,
		intro.IntroDate,
		intro.AgreedFee,
		intro.FeeType,
		intro.CreatedBy,
		intro.CreatedAt,
		intro.ModifiedBy,
		intro.LastModified,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *postgresIntroductionsRepo) GetByID(ctx context.Context, tenantID string, id int64) (*models.Introductions, error) {
	query := `
	SELECT id, tenant_id, introducer_id, introduced_party, property_id, transaction_id,
	       intro_date, agreed_fee, fee_type, created_by, created_at, modified_by, last_modified, deleted
	FROM introductions
	WHERE tenant_id = ? AND id = ? AND deleted = 0;
	`
	row := r.db.QueryRowContext(ctx, query, tenantID, id)

	var intro models.Introductions
	var deletedInt int
	err := row.Scan(
		&intro.ID,
		&intro.TenantID,
		&intro.IntroducerID,
		&intro.IntroducedParty,
		&intro.PropertyID,
		&intro.TransactionID,
		&intro.IntroDate,
		&intro.AgreedFee,
		&intro.FeeType,
		&intro.CreatedBy,
		&intro.CreatedAt,
		&intro.ModifiedBy,
		&intro.LastModified,
		&deletedInt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	intro.Deleted = (deletedInt != 0)
	return &intro, nil
}

func (r *postgresIntroductionsRepo) ListAll(ctx context.Context, tenantID string) ([]*models.Introductions, error) {
	query := `
	SELECT id, tenant_id, introducer_id, introduced_party, property_id, transaction_id,
	       intro_date, agreed_fee, fee_type, created_by, created_at, modified_by, last_modified, deleted
	FROM introductions
	WHERE tenant_id = ? AND deleted = 0;
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Introductions
	for rows.Next() {
		var intro models.Introductions
		var deletedInt int
		if err := rows.Scan(
			&intro.ID,
			&intro.TenantID,
			&intro.IntroducerID,
			&intro.IntroducedParty,
			&intro.PropertyID,
			&intro.TransactionID,
			&intro.IntroDate,
			&intro.AgreedFee,
			&intro.FeeType,
			&intro.CreatedBy,
			&intro.CreatedAt,
			&intro.ModifiedBy,
			&intro.LastModified,
			&deletedInt,
		); err != nil {
			return nil, err
		}
		intro.Deleted = (deletedInt != 0)
		out = append(out, &intro)
	}
	return out, nil
}

func (r *postgresIntroductionsRepo) Update(ctx context.Context, intro *models.Introductions) error {
	existing, err := r.GetByID(ctx, intro.TenantID, intro.ID)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	now := time.Now().UTC()
	intro.LastModified = now

	query := `
	UPDATE introductions
	SET introducer_id = ?, introduced_party = ?, property_id = ?, transaction_id = ?, intro_date = ?, agreed_fee = ?, fee_type = ?,
	    modified_by = ?, last_modified = ?, deleted = ?
	WHERE tenant_id = ? AND id = ?;
	`
	_, err = r.db.ExecContext(ctx, query,
		intro.IntroducerID,
		intro.IntroducedParty,
		intro.PropertyID,
		intro.TransactionID,
		intro.IntroDate,
		intro.AgreedFee,
		intro.FeeType,
		intro.ModifiedBy,
		intro.LastModified,
		boolToInt(intro.Deleted),
		intro.TenantID,
		intro.ID,
	)
	return err
}

func (r *postgresIntroductionsRepo) Delete(ctx context.Context, tenantID string, id int64) error {
	existing, err := r.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	query := `
	UPDATE introductions
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
