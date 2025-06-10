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

/ NewDBLocationRepo selects the concrete implementation based on driver.
func NewDBLocationRepo(db *sql.DB, driver string) LocationRepo {
	switch driver {
	case "postgres":
		return &postgresLocationRepo{db: db}
	case "oracle":
		return &oracleLocationRepo{db: db}
	case "sqlite":
		return &sqliteLocationRepo{db: db}
	default:
		panic("unsupported driver: " + driver)
	}
}
