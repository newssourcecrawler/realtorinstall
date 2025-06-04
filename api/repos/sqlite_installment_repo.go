package repos

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type sqliteInstallmentRepo struct{ db *sql.DB }

// NewSQLiteInstallmentRepo opens/creates "installments" table
func NewSQLiteInstallmentRepo(dbPath string) (InstallmentRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
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
	  last_modified DATETIME NOT NULL
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &sqliteInstallmentRepo{db: db}, nil
}

func (r *sqliteInstallmentRepo) Create(ctx context.Context, inst *models.Installment) (int64, error) {
	return 0, nil
}
func (r *sqliteInstallmentRepo) GetByID(ctx context.Context, id int64) (*models.Installment, error) {
	return nil, nil
}
func (r *sqliteInstallmentRepo) ListAll(ctx context.Context) ([]*models.Installment, error) {
	return nil, nil
}
func (r *sqliteInstallmentRepo) Update(ctx context.Context, inst *models.Installment) error {
	return nil
}
func (r *sqliteInstallmentRepo) Delete(ctx context.Context, id int64) error {
	return nil
}
