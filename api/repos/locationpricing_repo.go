package repos

import (
	"context"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

// LocationPricingRepo defines CRUD for LocationPricing
type LocationPricingRepo interface {
	Create(ctx context.Context, lp *models.LocationPricing) (int64, error)
	GetByID(ctx context.Context, tenantID string, id int64) (*models.LocationPricing, error)
	ListAll(ctx context.Context, tenantID string) ([]*models.LocationPricing, error)
	Update(ctx context.Context, lp *models.LocationPricing) error // use lp.TenantID
	Delete(ctx context.Context, tenantID string, id int64) error
}
