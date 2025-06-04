package repos

import (
	"context"

	"github.com/newssourcecrawler/realtorinstall/internal/models"
)

type BuyerRepo interface {
	Create(ctx context.Context, b *models.Buyer) (int64, error)
	GetByID(ctx context.Context, id string) (*models.Buyer, error)
	ListAll(ctx context.Context) ([]*models.Buyer, error)
	Update(ctx context.Context, id string, b *models.Buyer) error
	//Delete(ctx context.Context, id int64) error
}
