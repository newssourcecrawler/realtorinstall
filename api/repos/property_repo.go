package repos

import (
	"context"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

// PropertyRepo defines CRUD for Property
type PropertyRepo interface {
	Create(ctx context.Context, p *models.Property) (int64, error)
	GetByID(ctx context.Context, tenantID string, id int64) (*models.Property, error)
	ListAll(ctx context.Context, tenantID string) ([]*models.Property, error)
	Update(ctx context.Context, p *models.Property) error
	Delete(ctx context.Context, tenantID string, id int64) error
}
