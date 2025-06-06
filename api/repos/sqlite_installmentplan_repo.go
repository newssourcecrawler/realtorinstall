// api/repos/sqlite_installmentplan_repo.go

package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type sqliteInstallmentPlanRepo struct {
	db *sql.DB
}

func NewSQLiteInstallmentPlanRepo(dbPath string) (InstallmentPlanRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	schema := `
	CREATE TABLE IF NOT EXISTS installment_plans (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  tenant_id TEXT NOT NULL,
	  property_id INTEGER NOT NULL,
	  buyer_id INTEGER NOT NULL,
	  total_price REAL NOT NULL,
	  down_payment REAL NOT NULL,
	  num_installments INTEGER NOT NULL,
	  frequency TEXT NOT NULL,
	  first_installment DATETIME NOT NULL,
	  interest_rate REAL NOT NULL,
	  created_by TEXT NOT NULL,
	  created_at DATETIME NOT NULL,
	  modified_by TEXT NOT NULL,
	  last_modified DATETIME NOT NULL,
	  deleted INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_installmentplans_tenant ON installment_plans(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_installmentplans_property ON installment_plans(property_id);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &sqliteInstallmentPlanRepo{db: db}, nil
}

func (r *sqliteInstallmentPlanRepo) Create(ctx context.Context, p *models.InstallmentPlan) (int64, error) {
	if p.TenantID == "" || p.PropertyID == 0 || p.BuyerID == 0 || p.CreatedBy == "" || p.ModifiedBy == "" {
		return 0, errors.New("missing required fields or tenant/audit info")
	}
	now := time.Now().UTC()
	p.CreatedAt = now
	p.LastModified = now

	query := `
	INSERT INTO installment_plans (
	  tenant_id, property_id, buyer_id, total_price, down_payment, num_installments, frequency, first_installment, interest_rate,
	  created_by, created_at, modified_by, last_modified, deleted
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0);
	`
	res, err := r.db.ExecContext(ctx, query,
		p.TenantID,
		p.PropertyID,
		p.BuyerID,
		p.TotalPrice,
		p.DownPayment,
		p.NumInstallments,
		p.Frequency,
		p.FirstInstallment,
		p.InterestRate,
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

func (r *sqliteInstallmentPlanRepo) GetByID(ctx context.Context, tenantID string, id int64) (*models.InstallmentPlan, error) {
	query := `
	SELECT id, tenant_id, property_id, buyer_id, total_price, down_payment, num_installments, frequency, first_installment, interest_rate,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM installment_plans
	WHERE tenant_id = ? AND id = ? AND deleted = 0;
	`
	row := r.db.QueryRowContext(ctx, query, tenantID, id)

	var p models.InstallmentPlan
	var deletedInt int
	err := row.Scan(
		&p.ID,
		&p.TenantID,
		&p.PropertyID,
		&p.BuyerID,
		&p.TotalPrice,
		&p.DownPayment,
		&p.NumInstallments,
		&p.Frequency,
		&p.FirstInstallment,
		&p.InterestRate,
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

func (r *sqliteInstallmentPlanRepo) ListAll(ctx context.Context, tenantID string) ([]*models.InstallmentPlan, error) {
	query := `
	SELECT id, tenant_id, property_id, buyer_id, total_price, down_payment, num_installments, frequency, first_installment, interest_rate,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM installment_plans
	WHERE tenant_id = ? AND deleted = 0;
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.InstallmentPlan
	for rows.Next() {
		var p models.InstallmentPlan
		var deletedInt int
		if err := rows.Scan(
			&p.ID,
			&p.TenantID,
			&p.PropertyID,
			&p.BuyerID,
			&p.TotalPrice,
			&p.DownPayment,
			&p.NumInstallments,
			&p.Frequency,
			&p.FirstInstallment,
			&p.InterestRate,
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

func (r *sqliteInstallmentPlanRepo) ListByPlan(ctx context.Context, tenantID string, planID int64) ([]*models.InstallmentPlan, error) {
	query := `
	SELECT id, tenant_id, property_id, buyer_id, total_price, down_payment, num_installments, frequency, first_installment, interest_rate,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM installment_plans
	WHERE tenant_id = ? AND id = ? AND deleted = 0;
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.InstallmentPlan
	for rows.Next() {
		var p models.InstallmentPlan
		var deletedInt int
		if err := rows.Scan(
			&p.ID,
			&p.TenantID,
			&p.PropertyID,
			&p.BuyerID,
			&p.TotalPrice,
			&p.DownPayment,
			&p.NumInstallments,
			&p.Frequency,
			&p.FirstInstallment,
			&p.InterestRate,
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

func (r *sqliteInstallmentPlanRepo) Update(ctx context.Context, p *models.InstallmentPlan) error {
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
	UPDATE installment_plans
	SET property_id = ?, buyer_id = ?, total_price = ?, down_payment = ?, num_installments = ?, frequency = ?, first_installment = ?, interest_rate = ?,
	    modified_by = ?, last_modified = ?, deleted = ?
	WHERE tenant_id = ? AND id = ?;
	`
	_, err = r.db.ExecContext(ctx, query,
		p.PropertyID,
		p.BuyerID,
		p.TotalPrice,
		p.DownPayment,
		p.NumInstallments,
		p.Frequency,
		p.FirstInstallment,
		p.InterestRate,
		p.ModifiedBy,
		p.LastModified,
		boolToInt(p.Deleted),
		p.TenantID,
		p.ID,
	)
	return err
}

func (r *sqliteInstallmentPlanRepo) Delete(ctx context.Context, tenantID string, id int64) error {
	existing, err := r.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	query := `
	UPDATE installment_plans
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

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// PlanSummary holds plan‐ID and total outstanding balance.
type PlanSummary struct {
	PlanID           int64   `json:"plan_id"`
	TotalOutstanding float64 `json:"total_outstanding"`
}

// SummarizeByPlan computes “amount_due − amount_paid” grouped by each plan.
func (r *sqliteInstallmentPlanRepo) SummarizeByPlan(ctx context.Context, tenantID string) ([]models.PlanSummary, error) {
	query := `
        SELECT plan_id, 
            SUM(amount_due - amount_paid) AS total_outstanding
          FROM installments
         WHERE tenant_id = ? AND deleted = 0
         GROUP BY plan_id;
    `
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.PlanSummary
	for rows.Next() {
		var ps models.PlanSummary
		if err := rows.Scan(&ps.PlanID, &ps.TotalOutstanding); err != nil {
			return nil, err
		}
		out = append(out, ps)
	}
	return out, nil
}
