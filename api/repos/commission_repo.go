package repos

import (
	"context"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type CommissionRepo interface {
	Create(ctx context.Context, p *models.Commission) (int64, error) // p.TenantID set
	GetByID(ctx context.Context, tenantID string, id int64) (*models.Commission, error)
	ListAll(ctx context.Context, tenantID string) ([]*models.Commission, error)
	Update(ctx context.Context, b *models.Commission) error // using b.TenantID,b.ID
	Delete(ctx context.Context, tenantID string, id int64) error
	TotalCommissionByBeneficiary(ctx context.Context, tenantID string) ([]models.CommissionSummary, error)
	GetCommissionDetailsForBeneficiary(ctx context.Context, tenantID string, beneficiaryID int64) ([]*models.Commission, error)
}

/ NewDBCommissionRepo selects the concrete implementation based on driver.
func NewDBCommissionRepo(db *sql.DB, driver string) CommissionRepo {
	switch driver {
	case "postgres":
		return &postgresCommissionRepo{db: db}
	case "oracle":
		return &oracleCommissionRepo{db: db}
	case "sqlite":
		return &sqliteCommissionRepo{db: db}
	default:
		panic("unsupported driver: " + driver)
	}
}
