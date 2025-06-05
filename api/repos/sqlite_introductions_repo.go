package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/newssourcecrawler/realtorcommall/api/models"
)

type sqliteIntroductionsRepo struct {
	db *sql.DB
}

func NewSQLiteIntroductionsRepo(dbPath string) (IntroductionsRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	schema := `
	CREATE TABLE introductions (
		id                INTEGER PRIMARY KEY AUTOINCREMENT,
		tenant_id         TEXT    NOT NULL,     -- multi‐tenant support
		transaction_type  TEXT    NOT NULL,     -- “sale” | “letting” | “introduction”
		transaction_id    INTEGER NOT NULL,     -- e.g. sale.id, letting.id, or intro.id
		beneficiary_id    INTEGER NOT NULL,     -- foreign key to users/employees table, or an external party
		Introductions_type   TEXT    NOT NULL,     -- “percentage” | “fixed” | “credit”
		rate_or_amount    REAL    NOT NULL,     -- if percentage, store % as decimal (0.03); if fixed or credit, store flat amount
		calculated_amount REAL    NOT NULL,     -- actual money owed (pre‐tax/fees)
		memo              TEXT,                 -- optional description
		created_by        TEXT    NOT NULL,
		created_at        DATETIME NOT NULL,
		modified_by       TEXT    NOT NULL,
		last_modified     DATETIME NOT NULL,
		deleted           INTEGER NOT NULL DEFAULT 0
		);
	CREATE INDEX idx_introductions_tenant ON introductions(tenant_id);
	CREATE INDEX idx_introductions_txn     ON introductions(transaction_type, transaction_id);
	CREATE INDEX idx_introductions_benef    ON introductions(beneficiary_id);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &sqliteIntroductionsRepo{db: db}, nil
}

func (r *sqliteIntroductionsRepo) Create(ctx context.Context, comm *models.Introductions) (int64, error) {
	if comm.TenantID == "" || comm.PlanID == 0 || comm.CreatedBy == "" || comm.ModifiedBy == "" {
		return 0, errors.New("missing required fields or tenant/audit info")
	}
	now := time.Now().UTC()
	comm.CreatedAt = now
	comm.LastModified = now

	query := `
	INSERT INTO introductions (
	  tenant_id, plan_id, sequence_number, due_date, amount_due, amount_paid, status, late_fee, paid_date,
	  created_by, created_at, modified_by, last_modified, deleted
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0);
	`
	res, err := r.db.ExecContext(ctx, query,
		comm.TenantID,
		comm.PlanID,
		comm.SequenceNumber,
		comm.DueDate,
		comm.AmountDue,
		comm.AmountPaid,
		comm.Status,
		comm.LateFee,
		comm.PaidDate,
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

func (r *sqliteIntroductionsRepo) GetByID(ctx context.Context, tenantID string, id int64) (*models.Introductions, error) {
	query := `
	SELECT id, tenant_id, plan_id, sequence_number, due_date, amount_due, amount_paid, status, late_fee, paid_date,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM introductions
	WHERE tenant_id = ? AND id = ? AND deleted = 0;
	`
	row := r.db.QueryRowContext(ctx, query, tenantID, id)

	var comm models.Introductions
	var deletedInt int
	err := row.Scan(
		&comm.ID,
		&comm.TenantID,
		&comm.PlanID,
		&comm.SequenceNumber,
		&comm.DueDate,
		&comm.AmountDue,
		&comm.AmountPaid,
		&comm.Status,
		&comm.LateFee,
		&comm.PaidDate,
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
	comm.Deleted = deletedInt != 0
	return &comm, nil
}

func (r *sqliteIntroductionsRepo) ListAll(ctx context.Context, tenantID string) ([]*models.Introductions, error) {
	query := `
	SELECT id, tenant_id, plan_id, sequence_number, due_date, amount_due, amount_paid, status, late_fee, paid_date,
	       created_by, created_at, modified_by, last_modified, deleted
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
		var comm models.Introductions
		var deletedInt int
		if err := rows.Scan(
			&comm.ID,
			&comm.TenantID,
			&comm.PlanID,
			&comm.SequenceNumber,
			&comm.DueDate,
			&comm.AmountDue,
			&comm.AmountPaid,
			&comm.Status,
			&comm.LateFee,
			&comm.PaidDate,
			&comm.CreatedBy,
			&comm.CreatedAt,
			&comm.ModifiedBy,
			&comm.LastModified,
			&deletedInt,
		); err != nil {
			return nil, err
		}
		comm.Deleted = deletedInt != 0
		out = append(out, &comm)
	}
	return out, nil
}

func (r *sqliteIntroductionsRepo) ListByBeneficiary(ctx context.Context, tenantID string, BeneficiaryID int64) ([]*models.Introductions, error) {
	query := `
	SELECT beneficiary_id,
      SUM(calculated_amount) AS total_owed
		FROM introductions
		WHERE tenant_id = ? AND deleted = 0
		GROUP BY beneficiary_id;
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID, BeneficiaryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Introductions
	for rows.Next() {
		var comm models.Introductions
		var deletedInt int
		if err := rows.Scan(
			&comm.ID,
			&comm.TenantID,
			&comm.BeneficiaryID,
			&comm.DueDate,
			&comm.AmountDue,
			&comm.AmountPaid,
			&comm.Status,
			&comm.LateFee,
			&comm.PaidDate,
			&comm.CreatedBy,
			&comm.CreatedAt,
			&comm.ModifiedBy,
			&comm.LastModified,
			&deletedInt,
		); err != nil {
			return nil, err
		}
		comm.Deleted = deletedInt != 0
		out = append(out, &comm)
	}
	return out, nil
}

func (r *sqliteIntroductionsRepo) ListByTransaction(ctx context.Context, tenantID string, TransactionType string) ([]*models.Introductions, error) {
	query := `
	SELECT transaction_type,
      COUNT(*)           AS count_records,
      SUM(calculated_amount) AS sum_amount
		FROM introductions
		WHERE tenant_id = ? AND deleted = 0
		GROUP BY transaction_type;
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID, TransactionType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Introductions
	for rows.Next() {
		var comm models.Introductions
		var deletedInt int
		if err := rows.Scan(
			&comm.ID,
			&comm.TenantID,
			&comm.BeneficiaryID,
			&comm.DueDate,
			&comm.AmountDue,
			&comm.AmountPaid,
			&comm.Status,
			&comm.LateFee,
			&comm.PaidDate,
			&comm.CreatedBy,
			&comm.CreatedAt,
			&comm.ModifiedBy,
			&comm.LastModified,
			&deletedInt,
		); err != nil {
			return nil, err
		}
		comm.Deleted = deletedInt != 0
		out = append(out, &comm)
	}
	return out, nil
}

func (r *sqliteIntroductionsRepo) ListByMonth(ctx context.Context, tenantID string) ([]*models.Introductions, error) {
	query := `
	SELECT strftime('%Y-%m', created_at) AS year_month,
      SUM(calculated_amount) AS month_total
		FROM introductions
		WHERE tenant_id = ? AND deleted = 0
		GROUP BY year_month
		ORDER BY year_month DESC;
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Introductions
	for rows.Next() {
		var comm models.Introductions
		var deletedInt int
		if err := rows.Scan(
			&comm.ID,
			&comm.TenantID,
			&comm.BeneficiaryID,
			&comm.DueDate,
			&comm.AmountDue,
			&comm.AmountPaid,
			&comm.Status,
			&comm.LateFee,
			&comm.PaidDate,
			&comm.CreatedBy,
			&comm.CreatedAt,
			&comm.ModifiedBy,
			&comm.LastModified,
			&deletedInt,
		); err != nil {
			return nil, err
		}
		comm.Deleted = deletedInt != 0
		out = append(out, &comm)
	}
	return out, nil
}

func (r *sqliteIntroductionsRepo) Update(ctx context.Context, comm *models.Introductions) error {
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
	UPDATE introductions
	SET plan_id = ?, sequence_number = ?, due_date = ?, amount_due = ?, amount_paid = ?, status = ?, late_fee = ?, paid_date = ?,
	    modified_by = ?, last_modified = ?, deleted = ?
	WHERE tenant_id = ? AND id = ?;
	`
	_, err = r.db.ExecContext(ctx, query,
		comm.PlanID,
		comm.SequenceNumber,
		comm.DueDate,
		comm.AmountDue,
		comm.AmountPaid,
		comm.Status,
		comm.LateFee,
		comm.PaidDate,
		comm.ModifiedBy,
		comm.LastModified,
		boolToInt(comm.Deleted),
		comm.TenantID,
		comm.ID,
	)
	return err
}

func (r *sqliteIntroductionsRepo) Delete(ctx context.Context, tenantID string, id int64) error {
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
