package repos

import (
	"context"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type LettingsRepo interface {
	Create(ctx context.Context, p *models.Lettings) (int64, error) // p.TenantID set
	GetByID(ctx context.Context, tenantID string, id int64) (*models.Lettings, error)
	ListAll(ctx context.Context, tenantID string) ([]*models.Lettings, error)
	Update(ctx context.Context, b *models.Lettings) error // using b.TenantID,b.ID
	//Delete(ctx context.Context, tenantID string, id int64) error
	Delete(ctx context.Context, b *models.Lettings) error
	SummarizeRentRoll(ctx context.Context, tenantID string) ([]models.RentRoll, error)
}

// RentRoll holds a Lettings_id and total rent currently active.
type RentRoll struct {
	LettingsID int64   `json:"Lettings_id"`
	TotalRent  float64 `json:"total_rent"`
}

/ NewDBLettingsRepo selects the concrete implementation based on driver.
func NewDBLettingsRepo(db *sql.DB, driver string) LettingsRepo {
	switch driver {
	case "postgres":
		return &postgresLettingsRepo{db: db}
	case "oracle":
		return &oracleLettingsRepo{db: db}
	case "sqlite":
		return &sqliteLettingsRepo{db: db}
	default:
		panic("unsupported driver: " + driver)
	}
}
