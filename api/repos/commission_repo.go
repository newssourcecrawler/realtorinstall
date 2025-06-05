package repos

import (
	"context"
	"database/sql"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type CommissionRepo interface {
	Create(ctx context.Context, p *models.Commission) (int64, error) // p.TenantID set
	GetByID(ctx context.Context, tenantID string, id int64) (*models.Commission, error)
	ListAll(ctx context.Context, tenantID string) ([]*models.Commission, error)
	Update(ctx context.Context, b *models.Commission) error // using b.TenantID,b.ID
	Delete(ctx context.Context, tenantID string, id int64) error
	QueryCommissionByBeneficiary(ctx context.Context, tenantID string) (*sql.Rows, error)
}
