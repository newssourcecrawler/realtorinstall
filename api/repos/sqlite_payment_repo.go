package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/newssourcecrawler/realtorinstall/api/models"
	//"github.com/newssourcecrawler/realtorinstall/api/repos"
	//"github.com/newssourcecrawler/realtorinstall/api/repos"
)

// sqlitePaymentRepo implements PaymentRepo using SQLite.
type sqlitePaymentRepo struct {
	db *sql.DB
}

// NewSQLitePaymentRepo opens/creates "payments" table.
func NewSQLitePaymentRepo(dbPath string) (PaymentRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Add a soft‐delete flag
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
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &sqlitePaymentRepo{db: db}, nil
}

// Create inserts a new Payment record.
func (r *sqlitePaymentRepo) Create(ctx context.Context, p *models.Payment) (int64, error) {
	// Validate required fields
	if p.InstallmentID == 0 {
		return 0, errors.New("installment_id is required")
	}
	if p.AmountPaid <= 0 {
		return 0, errors.New("amount_paid must be positive")
	}
	if p.PaymentDate.IsZero() {
		return 0, errors.New("payment_date is required")
	}

	now := time.Now().UTC()
	p.CreatedAt = now
	p.LastModified = now

	query := `
	INSERT INTO payments (
	  installment_id,
	  tenant_id,
	  amount_paid,
	  payment_date,
	  payment_method,
	  transaction_ref,
	  created_at,
	  last_modified,
	  deleted
	) VALUES (?, ?, ?, ?, ?, ?, ?, 0);
	`
	res, err := r.db.ExecContext(ctx, query,
		p.InstallmentID,
		p.TenantID,
		p.AmountPaid,
		p.PaymentDate,
		p.PaymentMethod,
		p.TransactionRef,
		p.CreatedAt,
		p.LastModified,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetByID retrieves a Payment by its primary key (even if soft‐deleted).
func (r *sqlitePaymentRepo) GetByID(ctx context.Context, id int64) (*models.Payment, error) {
	query := `
	SELECT id, installment_id, amount_paid, payment_date, payment_method, transaction_ref, created_at, last_modified, deleted
	FROM payments
	WHERE id = ?;
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var p models.Payment
	var deletedInt int
	err := row.Scan(
		&p.ID,
		&p.InstallmentID,
		&p.AmountPaid,
		&p.PaymentDate,
		&p.PaymentMethod,
		&p.TransactionRef,
		&p.CreatedAt,
		&p.LastModified,
		&deletedInt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ListAll returns all non‐deleted Payments.
func (r *sqlitePaymentRepo) ListAll(ctx context.Context) ([]*models.Payment, error) {
	query := `
	SELECT id, tenant_id, installment_id, amount_paid, payment_date, payment_method, transaction_ref, created_at, last_modified, deleted
	FROM payments
	WHERE deleted = 0;
	`
	rows, err := r.db.QueryContext(ctx, query)
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
			&p.CreatedAt,
			&p.LastModified,
			&deletedInt,
		); err != nil {
			return nil, err
		}
		out = append(out, &p)
	}
	return out, nil
}

// Update modifies an existing Payment (including the Deleted flag).
func (r *sqlitePaymentRepo) Update(ctx context.Context, p *models.Payment) error {
	query := `
	UPDATE payments
	SET
	  installment_id = ?,
	  amount_paid = ?,
	  payment_date = ?,
	  payment_method = ?,
	  transaction_ref = ?,
	  last_modified = ?,
	  deleted = ?
	WHERE id = ?;
	`
	p.LastModified = time.Now().UTC()

	_, err := r.db.ExecContext(ctx, query,
		p.InstallmentID,
		p.AmountPaid,
		p.PaymentDate,
		p.PaymentMethod,
		p.TransactionRef,
		p.LastModified,
		boolToInt(p.Deleted),
		p.ID,
	)
	return err
}

// Delete soft‐deletes a Payment by setting deleted = 1.
func (r *sqlitePaymentRepo) Delete(ctx context.Context, id int64) error {
	// Confirm existence
	_, err := r.GetByID(ctx, id)
	if err == errors.New("record are required") {
		return errors.New("record are required")
	} else if err != nil {
		return err
	}

	query := `
	UPDATE payments
	SET deleted = 1, last_modified = ?
	WHERE id = ?;
	`
	_, err = r.db.ExecContext(ctx, query, time.Now().UTC(), id)
	return err
}
