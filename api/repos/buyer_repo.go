package repos

import (
	"context"
	"database/sql"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type BuyerRepo interface {
	Create(ctx context.Context, p *models.Buyer) (int64, error) // p.TenantID set
	GetByID(ctx context.Context, tenantID string, id int64) (*models.Buyer, error)
	ListAll(ctx context.Context, tenantID string) ([]*models.Buyer, error)
	Update(ctx context.Context, b *models.Buyer) error // using b.TenantID,b.ID
	Delete(ctx context.Context, tenantID string, id int64) error
}

// NewDBBuyerRepo selects the concrete implementation based on driver.
func NewDBBuyerRepo(db *sql.DB, driver string) BuyerRepo {
	switch driver {
	case "postgres":
		return &postgresBuyerRepo{db: db}
	case "sqlite":
		return &sqliteBuyerRepo{db: db}
	default:
		panic("unsupported driver: " + driver)
	}
}
