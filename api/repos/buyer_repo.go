package repos

import (
	"context"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type BuyerRepo interface {
	Create(ctx context.Context, p *models.Buyer) (int64, error) // p.TenantID set
	GetByID(ctx context.Context, tenantID string, id int64) (*models.Buyer, error)
	ListAll(ctx context.Context, tenantID string) ([]*models.Buyer, error)
	Update(ctx context.Context, b *models.Buyer) error // using b.TenantID,b.ID
	Delete(ctx context.Context, tenantID string, id int64) error
}
