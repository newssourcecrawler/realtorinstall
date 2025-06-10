package repos

import (
	"context"
	"database/sql"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

// PropertyRepo defines CRUD for Property
type PropertyRepo interface {
	Create(ctx context.Context, p *models.Property) (int64, error)
	GetByID(ctx context.Context, tenantID string, id int64) (*models.Property, error)
	ListAll(ctx context.Context, tenantID string) ([]*models.Property, error)
	Update(ctx context.Context, p *models.Property) error
	Delete(ctx context.Context, tenantID string, id int64) error
	SummarizeTopProperties(ctx context.Context, tenantID string) ([]models.PropertyPaymentVolume, error)
}

// PropertyPaymentVolume holds a property_id and total paid so far.
type PropertyPaymentVolume struct {
	PropertyID      int64   `json:"property_id"`
	TotalPaidAmount float64 `json:"total_paid"`
}

// NewDBPropertyRepo selects the concrete implementation based on driver.
func NewDBPropertyRepo(db *sql.DB, driver string) PropertyRepo {
	switch driver {
	case "postgres":
		return &postgresPropertyRepo{db: db}
	case "oracle":
		return &oraclePropertyRepo{db: db}
	case "sqlite":
		return &sqlitePropertyRepo{db: db}
	default:
		panic("unsupported driver: " + driver)
	}
}
