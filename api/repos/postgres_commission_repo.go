package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type postgresCommissionRepo struct {
	db *sql.DB
}

func NewPostgresCommissionRepo(db *sql.DB) CommissionRepo {
	return &postgresCommissionRepo{db: db}
}

func (r *postgresCommissionRepo) Create(ctx context.Context, comm *models.Commission) (int64, error) {
	if comm.TenantID == "" ||
		comm.TransactionType == "" ||
		comm.TransactionID == 0 ||
		comm.BeneficiaryID == 0 ||
		comm.CommissionType == "" ||
		comm.CreatedBy == "" ||
		comm.ModifiedBy == "" {
		return 0, errors.New("missing required fields or tenant/audit info")
	}
	now := time.Now().UTC()
	comm.CreatedAt = now
	comm.LastModified = now

	query := `
	INSERT INTO commissions (
	  tenant_id, transaction_type, transaction_id, beneficiary_id,
	  commission_type, rate_or_amount, calculated_amount, memo,
	  created_by, created_at, modified_by, last_modified, deleted
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0);
	`
	res, err := r.db.ExecContext(ctx, query,
		comm.TenantID,
		comm.TransactionType,
		comm.TransactionID,
		comm.BeneficiaryID,
		comm.CommissionType,
		comm.RateOrAmount,
		comm.CalculatedAmount,
		comm.Memo,
		comm.CreatedBy,
		comm.CreatedAt,
		comm.ModifiedBy,
		comm.LastModified,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *postgresCommissionRepo) GetByID(ctx context.Context, tenantID string, id int64) (*models.Commission, error) {
	query := `
	SELECT id, tenant_id, transaction_type, transaction_id, beneficiary_id,
	       commission_type, rate_or_amount, calculated_amount, memo,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM commissions
	WHERE tenant_id = ? AND id = ? AND deleted = 0;
	`
	row := r.db.QueryRowContext(ctx, query, tenantID, id)

	var comm models.Commission
	var deletedInt int
	err := row.Scan(
		&comm.ID,
		&comm.TenantID,
		&comm.TransactionType,
		&comm.TransactionID,
		&comm.BeneficiaryID,
		&comm.CommissionType,
		&comm.RateOrAmount,
		&comm.CalculatedAmount,
		&comm.Memo,
		&comm.CreatedBy,
		&comm.CreatedAt,
		&comm.ModifiedBy,
		&comm.LastModified,
		&deletedInt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	comm.Deleted = (deletedInt != 0)
	return &comm, nil
}

func (r *postgresCommissionRepo) ListAll(ctx context.Context, tenantID string) ([]*models.Commission, error) {
	query := `
	SELECT id, tenant_id, transaction_type, transaction_id, beneficiary_id,
	       commission_type, rate_or_amount, calculated_amount, memo,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM commissions
	WHERE tenant_id = ? AND deleted = 0;
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Commission
	for rows.Next() {
		var comm models.Commission
		var deletedInt int
		if err := rows.Scan(
			&comm.ID,
			&comm.TenantID,
			&comm.TransactionType,
			&comm.TransactionID,
			&comm.BeneficiaryID,
			&comm.CommissionType,
			&comm.RateOrAmount,
			&comm.CalculatedAmount,
			&comm.Memo,
			&comm.CreatedBy,
			&comm.CreatedAt,
			&comm.ModifiedBy,
			&comm.LastModified,
			&deletedInt,
		); err != nil {
			return nil, err
		}
		comm.Deleted = (deletedInt != 0)
		out = append(out, &comm)
	}
	return out, nil
}

func (r *postgresCommissionRepo) Update(ctx context.Context, comm *models.Commission) error {
	existing, err := r.GetByID(ctx, comm.TenantID, comm.ID)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	now := time.Now().UTC()
	comm.LastModified = now

	query := `
	UPDATE commissions
	SET transaction_type = ?, transaction_id = ?, beneficiary_id = ?,
	    commission_type = ?, rate_or_amount = ?, calculated_amount = ?, memo = ?,
	    modified_by = ?, last_modified = ?, deleted = ?
	WHERE tenant_id = ? AND id = ?;
	`
	_, err = r.db.ExecContext(ctx, query,
		comm.TransactionType,
		comm.TransactionID,
		comm.BeneficiaryID,
		comm.CommissionType,
		comm.RateOrAmount,
		comm.CalculatedAmount,
		comm.Memo,
		comm.ModifiedBy,
		comm.LastModified,
		boolToInt(comm.Deleted),
		comm.TenantID,
		comm.ID,
	)
	return err
}

func (r *postgresCommissionRepo) Delete(ctx context.Context, tenantID string, id int64) error {
	existing, err := r.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	query := `
	UPDATE commissions
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

// TotalCommissionByBeneficiary sums all earned commissions per beneficiary.
func (r *postgresCommissionRepo) TotalCommissionByBeneficiary(ctx context.Context, tenantID string) ([]models.CommissionSummary, error) {
	query := `
        SELECT beneficiary_id, SUM(calculated_amount) AS total_commission
          FROM commissions
         WHERE tenant_id = ? AND deleted = 0
         GROUP BY beneficiary_id;
    `
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.CommissionSummary
	for rows.Next() {
		var cs models.CommissionSummary
		if err := rows.Scan(&cs.BeneficiaryID, &cs.TotalCommission); err != nil {
			return nil, err
		}
		out = append(out, cs)
	}
	return out, nil
}

// TotalCommissionByBeneficiary sums all earned commissions per beneficiary.
func (r *postgresCommissionRepo) GetCommissionDetailsForBeneficiary(
	ctx context.Context,
	tenantID string,
	beneficiaryID int64,
) ([]*models.Commission, error) {

	query := `
        SELECT 
          id,
          tenant_id,
          transaction_type,
          transaction_id,
          beneficiary_id,
          commission_type,
          rate_or_amount,
          calculated_amount,
          memo,
          created_by,
          created_at,
          modified_by,
          last_modified,
          deleted
        FROM commissions
        WHERE tenant_id = ? 
          AND beneficiary_id = ? 
          AND deleted = 0;
    `
	rows, err := r.db.QueryContext(ctx, query, tenantID, beneficiaryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Commission
	for rows.Next() {
		var c models.Commission
		if err := rows.Scan(
			&c.ID,
			&c.TenantID,
			&c.TransactionType,
			&c.TransactionID,
			&c.BeneficiaryID,
			&c.CommissionType,
			&c.RateOrAmount,
			&c.CalculatedAmount,
			&c.Memo,
			&c.CreatedBy,
			&c.CreatedAt,
			&c.ModifiedBy,
			&c.LastModified,
			&c.Deleted,
		); err != nil {
			return nil, err
		}
		out = append(out, &c)
	}
	return out, rows.Err()
}
