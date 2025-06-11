package repos

import (
	"context"
	"database/sql"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

// PaymentRepo defines CRUD for Payment
type PaymentRepo interface {
	Create(ctx context.Context, p *models.Payment) (int64, error) // p.TenantID set
	GetByID(ctx context.Context, tenantID string, id int64) (*models.Payment, error)
	ListAll(ctx context.Context, tenantID string) ([]*models.Payment, error)
	ListByInstallment(ctx context.Context, tenantID string, installmentID int64) ([]*models.Payment, error)
	Update(ctx context.Context, p *models.Payment) error
	Delete(ctx context.Context, tenantID string, id int64) error
}

// NewDBPaymentRepo selects the concrete implementation based on driver.
func NewDBPaymentRepo(db *sql.DB, driver string) PaymentRepo {
	switch driver {
	case "postgres":
		return &postgresPaymentRepo{db: db}
	case "sqlite":
		return &sqlitePaymentRepo{db: db}
	default:
		panic("unsupported driver: " + driver)
	}
}
