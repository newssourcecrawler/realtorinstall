package repos

import (
	"context"

	"github.com/newssourcecrawler/realtorinstall/internal/models"
)

// LocationPricingRepo defines CRUD for LocationPricing
type LocationPricingRepo interface {
	Create(ctx context.Context, lp *models.LocationPricing) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.LocationPricing, error)
	ListAll(ctx context.Context) ([]*models.LocationPricing, error)
	Update(ctx context.Context, lp *models.LocationPricing) error
	Delete(ctx context.Context, id int64) error
}
