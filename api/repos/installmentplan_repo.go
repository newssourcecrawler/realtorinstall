package repos

import (
	"context"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

// InstallmentPlanRepo defines CRUD for InstallmentPlan
type InstallmentPlanRepo interface {
	Create(ctx context.Context, p *models.InstallmentPlan) (int64, error) // p.TenantID must be set
	GetByID(ctx context.Context, tenantID string, id int64) (*models.InstallmentPlan, error)
	ListAll(ctx context.Context, tenantID string) ([]*models.InstallmentPlan, error)
	ListByPlan(ctx context.Context, tenantID string, planID int64) ([]*models.InstallmentPlan, error)
	Update(ctx context.Context, p *models.InstallmentPlan) error // p.TenantID and p.ID must be set
	Delete(ctx context.Context, tenantID string, id int64) error
	SummarizeByPlan(ctx context.Context, tenantID string) ([]models.PlanSummary, error)
}

// PlanSummary holds plan‐ID and total outstanding balance.
type PlanSummary struct {
	PlanID           int64   `json:"plan_id"`
	TotalOutstanding float64 `json:"total_outstanding"`
}

/ NewDBInstallmentPlanRepo selects the concrete implementation based on driver.
func NewDBInstallmentPlanRepo(db *sql.DB, driver string) InstallmentPlanRepo {
	switch driver {
	case "postgres":
		return &postgresInstallmentPlanRepo{db: db}
	case "oracle":
		return &oracleInstallmentPlanRepo{db: db}
	case "sqlite":
		return &sqliteInstallmentPlanRepo{db: db}
	default:
		panic("unsupported driver: " + driver)
	}
}
