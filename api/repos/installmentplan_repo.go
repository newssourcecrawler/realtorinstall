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
	Update(ctx context.Context, p *models.InstallmentPlan) error
	Delete(ctx context.Context, tenantID string, id int64) error
}
