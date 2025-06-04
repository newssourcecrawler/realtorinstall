package repos

import (
	"context"

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
