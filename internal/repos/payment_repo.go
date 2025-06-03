package repos

import (
	"context"

	"github.com/newssourcecrawler/realtorinstall/internal/models"
)

// PaymentRepo defines CRUD for Payment
type PaymentRepo interface {
	Create(ctx context.Context, p *models.Payment) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.Payment, error)
	ListAll(ctx context.Context) ([]*models.Payment, error)
	Update(ctx context.Context, p *models.Payment) error
	Delete(ctx context.Context, id int64) error
}
