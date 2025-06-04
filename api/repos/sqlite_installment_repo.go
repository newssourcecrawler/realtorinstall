package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/newssourcecrawler/realtorinstall/api/models"
)

// sqliteInstallmentRepo implements InstallmentRepo using SQLite.
type sqliteInstallmentRepo struct {
	db *sql.DB
}

// NewSQLiteInstallmentRepo opens/creates "installments" table.
func NewSQLiteInstallmentRepo(dbPath string) (InstallmentRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Add a soft‐delete column
	schema := `
	CREATE TABLE IF NOT EXISTS installments (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  plan_id INTEGER NOT NULL,
	  sequence_number INTEGER NOT NULL,
	  due_date DATETIME NOT NULL,
	  amount_due REAL NOT NULL,
	  amount_paid REAL NOT NULL,
	  status TEXT NOT NULL,
	  late_fee REAL NOT NULL,
	  paid_date DATETIME,
	  created_at DATETIME NOT NULL,
	  last_modified DATETIME NOT NULL,
	  deleted INTEGER NOT NULL DEFAULT 0
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &sqliteInstallmentRepo{db: db}, nil
}

// Create inserts a new Installment.
func (r *sqliteInstallmentRepo) Create(ctx context.Context, inst *models.Installment) (int64, error) {
	// Validate required fields
	if inst.PlanID == 0 {
		return 0, errors.New("plan_id is required")
	}
	if inst.SequenceNumber <= 0 {
		return 0, errors.New("sequence_number must be positive")
	}
	if inst.DueDate.IsZero() {
		return 0, errors.New("due_date is required")
	}
	if inst.Status == "" {
		return 0, errors.New("status is required")
	}

	now := time.Now().UTC()
	inst.CreatedAt = now
	inst.LastModified = now

	query := `
	INSERT INTO installments (
	  plan_id,
	  sequence_number,
	  due_date,
	  amount_due,
	  amount_paid,
	  status,
	  late_fee,
	  paid_date,
	  created_at,
	  last_modified,
	  deleted
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0);
	`
	res, err := r.db.ExecContext(ctx, query,
		inst.PlanID,
		inst.SequenceNumber,
		inst.DueDate,
		inst.AmountDue,
		inst.AmountPaid,
		inst.Status,
		inst.LateFee,
		nullTime(inst.PaidDate),
		inst.CreatedAt,
		inst.LastModified,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetByID retrieves one Installment by its primary key.
func (r *sqliteInstallmentRepo) GetByID(ctx context.Context, id int64) (*models.Installment, error) {
	query := `
	SELECT id, plan_id, sequence_number, due_date, amount_due, amount_paid, status, late_fee, paid_date, created_at, last_modified, deleted
	FROM installments
	WHERE id = ?;
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var inst models.Installment
	var paidDate sql.NullTime
	var deletedInt int
	err := row.Scan(
		&inst.ID,
		&inst.PlanID,
		&inst.SequenceNumber,
		&inst.DueDate,
		&inst.AmountDue,
		&inst.AmountPaid,
		&inst.Status,
		&inst.LateFee,
		&paidDate,
		&inst.CreatedAt,
		&inst.LastModified,
		&deletedInt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if paidDate.Valid {
		inst.PaidDate = paidDate.Time
	}
	return &inst, nil
}

// ListAll returns all non‐deleted Installments.
func (r *sqliteInstallmentRepo) ListAll(ctx context.Context) ([]*models.Installment, error) {
	query := `
	SELECT id, plan_id, sequence_number, due_date, amount_due, amount_paid, status, late_fee, paid_date, created_at, last_modified, deleted
	FROM installments
	WHERE deleted = 0;
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Installment
	for rows.Next() {
		var inst models.Installment
		var paidDate sql.NullTime
		var deletedInt int
		if err := rows.Scan(
			&inst.ID,
			&inst.PlanID,
			&inst.SequenceNumber,
			&inst.DueDate,
			&inst.AmountDue,
			&inst.AmountPaid,
			&inst.Status,
			&inst.LateFee,
			&paidDate,
			&inst.CreatedAt,
			&inst.LastModified,
			&deletedInt,
		); err != nil {
			return nil, err
		}
		if paidDate.Valid {
			inst.PaidDate = paidDate.Time
		}
		out = append(out, &inst)
	}
	return out, nil
}

// Update modifies an existing Installment (including the deleted flag).
func (r *sqliteInstallmentRepo) Update(ctx context.Context, inst *models.Installment) error {
	query := `
	UPDATE installments
	SET
	  plan_id = ?,
	  sequence_number = ?,
	  due_date = ?,
	  amount_due = ?,
	  amount_paid = ?,
	  status = ?,
	  late_fee = ?,
	  paid_date = ?,
	  last_modified = ?,
	  deleted = ?
	WHERE id = ?;
	`
	inst.LastModified = time.Now().UTC()
	_, err := r.db.ExecContext(ctx, query,
		inst.PlanID,
		inst.SequenceNumber,
		inst.DueDate,
		inst.AmountDue,
		inst.AmountPaid,
		inst.Status,
		inst.LateFee,
		nullTime(inst.PaidDate),
		inst.LastModified,
		boolToInt(inst.Deleted),
		inst.ID,
	)
	return err
}

// Delete soft‐deletes an Installment.
func (r *sqliteInstallmentRepo) Delete(ctx context.Context, id int64) error {
	// Confirm it exists
	_, err := r.GetByID(ctx, id)
	if err == ErrNotFound {
		return ErrNotFound
	} else if err != nil {
		return err
	}

	query := `
	UPDATE installments
	SET deleted = 1, last_modified = ?
	WHERE id = ?;
	`
	_, err = r.db.ExecContext(ctx, query, time.Now().UTC(), id)
	return err
}

// Helper: convert nil time.Time to sql.NullTime
func nullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}
