package repos

import (
	"context"
	"database/sql"

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

// NewDBLocationRepo selects the concrete implementation based on driver.
func NewDBLocationRepo(db *sql.DB, driver string) LocationPricingRepo {
	switch driver {
	case "postgres":
		return &postgresLocationPricingRepo{db: db}
	case "sqlite":
		return &sqliteLocationPricingRepo{db: db}
	default:
		panic("unsupported driver: " + driver)
	}
}
