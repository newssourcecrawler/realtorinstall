package repos

import (
	"context"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

// InstallmentRepo defines CRUD for Installment
type InstallmentRepo interface {
	Create(ctx context.Context, inst *models.Installment) (int64, error) // inst.TenantID must be set
	GetByID(ctx context.Context, tenantID string, id int64) (*models.Installment, error)
	ListAll(ctx context.Context, tenantID string) ([]*models.Installment, error)
	ListByPlan(ctx context.Context, tenantID string, planID int64) ([]*models.Installment, error)
	Update(ctx context.Context, inst *models.Installment) error
	Delete(ctx context.Context, tenantID string, id int64) error
}
