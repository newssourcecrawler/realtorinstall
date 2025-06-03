package repos

import (
	"context"

	"github.com/newssourcecrawler/realtorinstall/internal/models"
)

// InstallmentPlanRepo defines CRUD for InstallmentPlan
type InstallmentPlanRepo interface {
	Create(ctx context.Context, p *models.InstallmentPlan) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.InstallmentPlan, error)
	ListAll(ctx context.Context) ([]*models.InstallmentPlan, error)
	Update(ctx context.Context, p *models.InstallmentPlan) error
	Delete(ctx context.Context, id int64) error
}
