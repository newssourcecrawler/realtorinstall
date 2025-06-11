package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type postgresInstallmentRepo struct {
	db *sql.DB
}

func NewPostgresInstallmentRepo(db *sql.DB) InstallmentRepo {
	return &postgresInstallmentRepo{db: db}
}

func (r *postgresInstallmentRepo) Create(ctx context.Context, inst *models.Installment) (int64, error) {
	if inst.TenantID == "" || inst.PlanID == 0 || inst.CreatedBy == "" || inst.ModifiedBy == "" {
		return 0, errors.New("missing required fields or tenant/audit info")
	}
	now := time.Now().UTC()
	inst.CreatedAt = now
	inst.LastModified = now

	query := `
	INSERT INTO installments (
	  tenant_id, plan_id, sequence_number, due_date, amount_due, amount_paid, status, late_fee, paid_date,
	  created_by, created_at, modified_by, last_modified, deleted
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0);
	`
	res, err := r.db.ExecContext(ctx, query,
		inst.TenantID,
		inst.PlanID,
		inst.SequenceNumber,
		inst.DueDate,
		inst.AmountDue,
		inst.AmountPaid,
		inst.Status,
		inst.LateFee,
		inst.PaidDate,
		inst.CreatedBy,
		inst.CreatedAt,
		inst.ModifiedBy,
		inst.LastModified,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *postgresInstallmentRepo) GetByID(ctx context.Context, tenantID string, id int64) (*models.Installment, error) {
	query := `
	SELECT id, tenant_id, plan_id, sequence_number, due_date, amount_due, amount_paid, status, late_fee, paid_date,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM installments
	WHERE tenant_id = ? AND id = ? AND deleted = 0;
	`
	row := r.db.QueryRowContext(ctx, query, tenantID, id)

	var inst models.Installment
	var deletedInt int
	err := row.Scan(
		&inst.ID,
		&inst.TenantID,
		&inst.PlanID,
		&inst.SequenceNumber,
		&inst.DueDate,
		&inst.AmountDue,
		&inst.AmountPaid,
		&inst.Status,
		&inst.LateFee,
		&inst.PaidDate,
		&inst.CreatedBy,
		&inst.CreatedAt,
		&inst.ModifiedBy,
		&inst.LastModified,
		&deletedInt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	inst.Deleted = deletedInt != 0
	return &inst, nil
}

func (r *postgresInstallmentRepo) ListAll(ctx context.Context, tenantID string) ([]*models.Installment, error) {
	query := `
	SELECT id, tenant_id, plan_id, sequence_number, due_date, amount_due, amount_paid, status, late_fee, paid_date,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM installments
	WHERE tenant_id = ? AND deleted = 0;
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Installment
	for rows.Next() {
		var inst models.Installment
		var deletedInt int
		if err := rows.Scan(
			&inst.ID,
			&inst.TenantID,
			&inst.PlanID,
			&inst.SequenceNumber,
			&inst.DueDate,
			&inst.AmountDue,
			&inst.AmountPaid,
			&inst.Status,
			&inst.LateFee,
			&inst.PaidDate,
			&inst.CreatedBy,
			&inst.CreatedAt,
			&inst.ModifiedBy,
			&inst.LastModified,
			&deletedInt,
		); err != nil {
			return nil, err
		}
		inst.Deleted = deletedInt != 0
		out = append(out, &inst)
	}
	return out, nil
}

func (r *postgresInstallmentRepo) ListByPlan(ctx context.Context, tenantID string, planID int64) ([]*models.Installment, error) {
	query := `
	SELECT id, tenant_id, plan_id, sequence_number, due_date, amount_due, amount_paid, status, late_fee, paid_date,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM installments
	WHERE tenant_id = ? AND plan_id = ? AND deleted = 0;
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Installment
	for rows.Next() {
		var inst models.Installment
		var deletedInt int
		if err := rows.Scan(
			&inst.ID,
			&inst.TenantID,
			&inst.PlanID,
			&inst.SequenceNumber,
			&inst.DueDate,
			&inst.AmountDue,
			&inst.AmountPaid,
			&inst.Status,
			&inst.LateFee,
			&inst.PaidDate,
			&inst.CreatedBy,
			&inst.CreatedAt,
			&inst.ModifiedBy,
			&inst.LastModified,
			&deletedInt,
		); err != nil {
			return nil, err
		}
		inst.Deleted = deletedInt != 0
		out = append(out, &inst)
	}
	return out, nil
}

func (r *postgresInstallmentRepo) Update(ctx context.Context, inst *models.Installment) error {
	existing, err := r.GetByID(ctx, inst.TenantID, inst.ID)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	now := time.Now().UTC()
	inst.LastModified = now

	query := `
	UPDATE installments
	SET plan_id = ?, sequence_number = ?, due_date = ?, amount_due = ?, amount_paid = ?, status = ?, late_fee = ?, paid_date = ?,
	    modified_by = ?, last_modified = ?, deleted = ?
	WHERE tenant_id = ? AND id = ?;
	`
	_, err = r.db.ExecContext(ctx, query,
		inst.PlanID,
		inst.SequenceNumber,
		inst.DueDate,
		inst.AmountDue,
		inst.AmountPaid,
		inst.Status,
		inst.LateFee,
		inst.PaidDate,
		inst.ModifiedBy,
		inst.LastModified,
		boolToInt(inst.Deleted),
		inst.TenantID,
		inst.ID,
	)
	return err
}

func (r *postgresInstallmentRepo) Delete(ctx context.Context, tenantID string, id int64) error {
	existing, err := r.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	query := `
	UPDATE installments
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
