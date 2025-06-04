package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/newssourcecrawler/realtorinstall/api/models"
)

// sqliteInstallmentPlanRepo implements InstallmentPlanRepo using SQLite.
type sqliteInstallmentPlanRepo struct {
	db *sql.DB
}

// NewSQLitePlanRepo opens/creates the "installment_plans" table.
func NewSQLitePlanRepo(dbPath string) (InstallmentPlanRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Add a 'deleted' column for soft‐delete
	schema := `
	CREATE TABLE IF NOT EXISTS installment_plans (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  property_id INTEGER NOT NULL,
	  buyer_id INTEGER NOT NULL,
	  total_price REAL NOT NULL,
	  down_payment REAL NOT NULL,
	  num_installments INTEGER NOT NULL,
	  frequency TEXT NOT NULL,
	  first_installment DATETIME NOT NULL,
	  interest_rate REAL NOT NULL,
	  created_at DATETIME NOT NULL,
	  last_modified DATETIME NOT NULL,
	  deleted INTEGER NOT NULL DEFAULT 0
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &sqliteInstallmentPlanRepo{db: db}, nil
}

// Create inserts a new InstallmentPlan.
func (r *sqliteInstallmentPlanRepo) Create(ctx context.Context, p *models.InstallmentPlan) (int64, error) {
	// Validate required fields
	if p.PropertyID == 0 {
		return 0, errors.New("property_id is required")
	}
	if p.BuyerID == 0 {
		return 0, errors.New("buyer_id is required")
	}
	if p.TotalPrice <= 0 {
		return 0, errors.New("total_price must be positive")
	}
	if p.NumInstallments <= 0 {
		return 0, errors.New("num_installments must be positive")
	}
	if p.Frequency == "" {
		return 0, errors.New("frequency is required")
	}
	if p.FirstInstallment.IsZero() {
		return 0, errors.New("first_installment is required")
	}

	now := time.Now().UTC()
	p.CreatedAt = now
	p.LastModified = now

	query := `
	INSERT INTO installment_plans (
	  property_id,
	  buyer_id,
	  total_price,
	  down_payment,
	  num_installments,
	  frequency,
	  first_installment,
	  interest_rate,
	  created_at,
	  last_modified,
	  deleted
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0);
	`
	res, err := r.db.ExecContext(ctx, query,
		p.PropertyID,
		p.BuyerID,
		p.TotalPrice,
		p.DownPayment,
		p.NumInstallments,
		p.Frequency,
		p.FirstInstallment,
		p.InterestRate,
		p.CreatedAt,
		p.LastModified,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetByID retrieves one InstallmentPlan by its primary key.
func (r *sqliteInstallmentPlanRepo) GetByID(ctx context.Context, id int64) (*models.InstallmentPlan, error) {
	query := `
	SELECT id, property_id, buyer_id, total_price, down_payment, num_installments, frequency, first_installment, interest_rate, created_at, last_modified, deleted
	FROM installment_plans
	WHERE id = ?;
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var p models.InstallmentPlan
	var deletedInt int
	err := row.Scan(
		&p.ID,
		&p.PropertyID,
		&p.BuyerID,
		&p.TotalPrice,
		&p.DownPayment,
		&p.NumInstallments,
		&p.Frequency,
		&p.FirstInstallment,
		&p.InterestRate,
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

// ListAll returns all non‐deleted InstallmentPlans.
func (r *sqliteInstallmentPlanRepo) ListAll(ctx context.Context) ([]*models.InstallmentPlan, error) {
	query := `
	SELECT id, property_id, buyer_id, total_price, down_payment, num_installments, frequency, first_installment, interest_rate, created_at, last_modified, deleted
	FROM installment_plans
	WHERE deleted = 0;
	`
	rows, err := r.db.QueryContext(ctx, query)
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
			&p.PropertyID,
			&p.BuyerID,
			&p.TotalPrice,
			&p.DownPayment,
			&p.NumInstallments,
			&p.Frequency,
			&p.FirstInstallment,
			&p.InterestRate,
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

// Update modifies an existing InstallmentPlan (including its deleted flag).
func (r *sqliteInstallmentPlanRepo) Update(ctx context.Context, p *models.InstallmentPlan) error {
	query := `
	UPDATE installment_plans
	SET
	  property_id = ?,
	  buyer_id = ?,
	  total_price = ?,
	  down_payment = ?,
	  num_installments = ?,
	  frequency = ?,
	  first_installment = ?,
	  interest_rate = ?,
	  last_modified = ?,
	  deleted = ?
	WHERE id = ?;
	`
	p.LastModified = time.Now().UTC()
	_, err := r.db.ExecContext(ctx, query,
		p.PropertyID,
		p.BuyerID,
		p.TotalPrice,
		p.DownPayment,
		p.NumInstallments,
		p.Frequency,
		p.FirstInstallment,
		p.InterestRate,
		p.LastModified,
		boolToInt(p.Deleted),
		p.ID,
	)
	return err
}

// Delete soft‐deletes an InstallmentPlan.
func (r *sqliteInstallmentPlanRepo) Delete(ctx context.Context, id int64) error {
	// Confirm existence
	_, err := r.GetByID(ctx, id)
	if err == ErrNotFound {
		return ErrNotFound
	} else if err != nil {
		return err
	}

	query := `
	UPDATE installment_plans
	SET deleted = 1, last_modified = ?
	WHERE id = ?;
	`
	_, err = r.db.ExecContext(ctx, query, time.Now().UTC(), id)
	return err
}

// Helper: convert bool → int
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
