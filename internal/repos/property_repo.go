// internal/repos/property_repo.go
package repos

import (
	"context"

	"github.com/newssourcecrawler/realtorinstall/internal/models"
)

// PropertyRepo defines CRUD for Property
type PropertyRepo interface {
	Create(ctx context.Context, p *models.Property) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.Property, error)
	ListAll(ctx context.Context) ([]*models.Property, error)
	Update(ctx context.Context, p *models.Property) error
	Delete(ctx context.Context, id int64) error
}
