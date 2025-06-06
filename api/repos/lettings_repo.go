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

// RentRoll holds a property_id and total rent currently active.
type RentRoll struct {
	PropertyID int64   `json:"property_id"`
	TotalRent  float64 `json:"total_rent"`
}
