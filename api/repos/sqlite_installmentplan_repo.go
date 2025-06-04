package repos

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // ensure SQLite driver is imported
	"github.com/newssourcecrawler/realtorinstall/api/models"
)

// InstallmentPlanRepo is already defined in internal/repos/installmentplan_repo.go
type sqliteInstallmentPlanRepo struct{ db *sql.DB }

// NewSQLitePlanRepo opens/creates the DB and ensures the "installment_plans" table exists.
func NewSQLitePlanRepo(dbPath string) (InstallmentPlanRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
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
	  last_modified DATETIME NOT NULL
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &sqliteInstallmentPlanRepo{db: db}, nil
}

// Stub methods to satisfy the interface; fill in later.
func (r *sqliteInstallmentPlanRepo) Create(ctx context.Context, p *models.InstallmentPlan) (int64, error) {
	return 0, nil
}
func (r *sqliteInstallmentPlanRepo) GetByID(ctx context.Context, id int64) (*models.InstallmentPlan, error) {
	return nil, nil
}
func (r *sqliteInstallmentPlanRepo) ListAll(ctx context.Context) ([]*models.InstallmentPlan, error) {
	return nil, nil
}
func (r *sqliteInstallmentPlanRepo) Update(ctx context.Context, p *models.InstallmentPlan) error {
	return nil
}
func (r *sqliteInstallmentPlanRepo) Delete(ctx context.Context, id int64) error {
	return nil
}
