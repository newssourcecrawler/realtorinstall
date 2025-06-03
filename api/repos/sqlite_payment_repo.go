package repos

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/newssourcecrawler/realtorinstall/internal/models"
)

type sqlitePaymentRepo struct{ db *sql.DB }

// NewSQLitePaymentRepo opens/creates "payments" table
func NewSQLitePaymentRepo(dbPath string) (PaymentRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	schema := `
	CREATE TABLE IF NOT EXISTS payments (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  installment_id INTEGER NOT NULL,
	  amount_paid REAL NOT NULL,
	  payment_date DATETIME NOT NULL,
	  payment_method TEXT NOT NULL,
	  transaction_ref TEXT,
	  created_at DATETIME NOT NULL,
	  last_modified DATETIME NOT NULL
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &sqlitePaymentRepo{db: db}, nil
}

func (r *sqlitePaymentRepo) Create(ctx context.Context, p *models.Payment) (int64, error) {
	return 0, nil
}
func (r *sqlitePaymentRepo) GetByID(ctx context.Context, id int64) (*models.Payment, error) {
	return nil, nil
}
func (r *sqlitePaymentRepo) ListAll(ctx context.Context) ([]*models.Payment, error) {
	return nil, nil
}
func (r *sqlitePaymentRepo) Update(ctx context.Context, p *models.Payment) error {
	return nil
}
func (r *sqlitePaymentRepo) Delete(ctx context.Context, id int64) error {
	return nil
}
