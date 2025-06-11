package repos

import (
	"context"
	"database/sql"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

// InstallmentRepo defines CRUD for Installment
type InstallmentRepo interface {
	Create(ctx context.Context, inst *models.Installment) (int64, error) // inst.TenantID must be set
	GetByID(ctx context.Context, tenantID string, id int64) (*models.Installment, error)
	ListAll(ctx context.Context, tenantID string) ([]*models.Installment, error)
	ListByPlan(ctx context.Context, tenantID string, planID int64) ([]*models.Installment, error)
	Update(ctx context.Context, inst *models.Installment) error // inst.TenantID and inst.ID must be set
	Delete(ctx context.Context, tenantID string, id int64) error
}

// NewDBInstallmentRepo selects the concrete implementation based on driver.
func NewDBInstallmentRepo(db *sql.DB, driver string) InstallmentRepo {
	switch driver {
	case "postgres":
		return &postgresInstallmentRepo{db: db}
	case "oracle":
		return &oracleInstallmentRepo{db: db}
	case "sqlite":
		return &sqliteInstallmentRepo{db: db}
	default:
		panic("unsupported driver: " + driver)
	}
}
