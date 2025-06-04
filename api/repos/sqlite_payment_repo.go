package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type sqlitePaymentRepo struct {
	db *sql.DB
}

func NewSQLitePaymentRepo(dbPath string) (PaymentRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	schema := `
	CREATE TABLE IF NOT EXISTS payments (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  tenant_id TEXT NOT NULL,
	  installment_id INTEGER NOT NULL,
	  amount_paid REAL NOT NULL,
	  payment_date DATETIME NOT NULL,
	  payment_method TEXT NOT NULL,
	  transaction_ref TEXT,
	  created_by TEXT NOT NULL,
	  created_at DATETIME NOT NULL,
	  modified_by TEXT NOT NULL,
	  last_modified DATETIME NOT NULL,
	  deleted INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_payments_tenant ON payments(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_payments_installment ON payments(installment_id);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &sqlitePaymentRepo{db: db}, nil
}

func (r *sqlitePaymentRepo) Create(ctx context.Context, p *models.Payment) (int64, error) {
	if p.TenantID == "" || p.InstallmentID == 0 || p.AmountPaid <= 0 || p.PaymentMethod == "" || p.CreatedBy == "" || p.ModifiedBy == "" {
		return 0, errors.New("missing required fields or tenant/audit info")
	}
	now := time.Now().UTC()
	p.CreatedAt = now
	p.LastModified = now
	query := `
	INSERT INTO payments (
	  tenant_id, installment_id, amount_paid, payment_date, payment_method, transaction_ref,
	  created_by, created_at, modified_by, last_modified, deleted
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0);
	`
	res, err := r.db.ExecContext(ctx, query,
		p.TenantID,
		p.InstallmentID,
		p.AmountPaid,
		p.PaymentDate,
		p.PaymentMethod,
		p.TransactionRef,
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

func (r *sqlitePaymentRepo) GetByID(ctx context.Context, tenantID string, id int64) (*models.Payment, error) {
	query := `
	SELECT id, tenant_id, installment_id, amount_paid, payment_date, payment_method, transaction_ref,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM payments
	WHERE tenant_id = ? AND id = ? AND deleted = 0;
	`
	row := r.db.QueryRowContext(ctx, query, tenantID, id)
	var p models.Payment
	var deletedInt int
	err := row.Scan(
		&p.ID,
		&p.TenantID,
		&p.InstallmentID,
		&p.AmountPaid,
		&p.PaymentDate,
		&p.PaymentMethod,
		&p.TransactionRef,
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

func (r *sqlitePaymentRepo) ListAll(ctx context.Context, tenantID string) ([]*models.Payment, error) {
	query := `
	SELECT id, tenant_id, installment_id, amount_paid, payment_date, payment_method, transaction_ref,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM payments
	WHERE tenant_id = ? AND deleted = 0;
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.Payment
	for rows.Next() {
		var p models.Payment
		var deletedInt int
		if err := rows.Scan(
			&p.ID,
			&p.TenantID,
			&p.InstallmentID,
			&p.AmountPaid,
			&p.PaymentDate,
			&p.PaymentMethod,
			&p.TransactionRef,
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

func (r *sqlitePaymentRepo) Update(ctx context.Context, p *models.Payment) error {
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
	UPDATE payments
	SET installment_id = ?, amount_paid = ?, payment_date = ?, payment_method = ?, transaction_ref = ?,
	    modified_by = ?, last_modified = ?, deleted = ?
	WHERE tenant_id = ? AND id = ?;
	`
	_, err = r.db.ExecContext(ctx, query,
		p.InstallmentID,
		p.AmountPaid,
		p.PaymentDate,
		p.PaymentMethod,
		p.TransactionRef,
		p.ModifiedBy,
		p.LastModified,
		boolToInt(p.Deleted),
		p.TenantID,
		p.ID,
	)
	return err
}

func (r *sqlitePaymentRepo) Delete(ctx context.Context, tenantID string, id int64) error {
	existing, err := r.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	query := `
	UPDATE payments
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
