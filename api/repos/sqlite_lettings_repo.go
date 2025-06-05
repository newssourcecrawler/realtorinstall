package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type sqliteLettingsRepo struct {
	db *sql.DB
}

func NewSQLiteLettingsRepo(dbPath string) (LettingsRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	schema := `
	CREATE TABLE IF NOT EXISTS lettings (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  tenant_id TEXT NOT NULL,
	  plan_id INTEGER NOT NULL,
	  sequence_number INTEGER NOT NULL,
	  due_date DATETIME NOT NULL,
	  amount_due REAL NOT NULL,
	  amount_paid REAL NOT NULL,
	  status TEXT NOT NULL,
	  late_fee REAL NOT NULL,
	  paid_date DATETIME,
	  created_by TEXT NOT NULL,
	  created_at DATETIME NOT NULL,
	  modified_by TEXT NOT NULL,
	  last_modified DATETIME NOT NULL,
	  deleted INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_Lettings_tenant ON Lettings(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_Lettings_plan ON Lettings(plan_id);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &sqliteLettingsRepo{db: db}, nil
}

func (r *sqliteLettingsRepo) Create(ctx context.Context, inst *models.Lettings) (int64, error) {
	if inst.TenantID == "" || inst.PlanID == 0 || inst.CreatedBy == "" || inst.ModifiedBy == "" {
		return 0, errors.New("missing required fields or tenant/audit info")
	}
	now := time.Now().UTC()
	inst.CreatedAt = now
	inst.LastModified = now

	query := `
	INSERT INTO lettings (
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

func (r *sqliteLettingsRepo) GetByID(ctx context.Context, tenantID string, id int64) (*models.Lettings, error) {
	query := `
	SELECT id, tenant_id, plan_id, sequence_number, due_date, amount_due, amount_paid, status, late_fee, paid_date,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM lettings
	WHERE tenant_id = ? AND id = ? AND deleted = 0;
	`
	row := r.db.QueryRowContext(ctx, query, tenantID, id)

	var inst models.Lettings
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

func (r *sqliteLettingsRepo) ListAll(ctx context.Context, tenantID string) ([]*models.Lettings, error) {
	query := `
	SELECT id, tenant_id, plan_id, sequence_number, due_date, amount_due, amount_paid, status, late_fee, paid_date,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM lettings
	WHERE tenant_id = ? AND deleted = 0;
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Lettings
	for rows.Next() {
		var inst models.Lettings
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

func (r *sqliteLettingsRepo) ListByPlan(ctx context.Context, tenantID string, planID int64) ([]*models.Lettings, error) {
	query := `
	SELECT id, tenant_id, plan_id, sequence_number, due_date, amount_due, amount_paid, status, late_fee, paid_date,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM lettings
	WHERE tenant_id = ? AND plan_id = ? AND deleted = 0;
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Lettings
	for rows.Next() {
		var inst models.Lettings
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

func (r *sqliteLettingsRepo) Update(ctx context.Context, inst *models.Lettings) error {
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
	UPDATE lettings
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

func (r *sqliteLettingsRepo) Delete(ctx context.Context, tenantID string, id int64) error {
	existing, err := r.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	query := `
	UPDATE lettings
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
