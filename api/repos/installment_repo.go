package repos

import (
	"context"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

// InstallmentRepo defines CRUD for Installment
type InstallmentRepo interface {
	Create(ctx context.Context, inst *models.Installment) (int64, error)
	GetByID(ctx context.Context, id int64) (*models.Installment, error)
	ListAll(ctx context.Context) ([]*models.Installment, error)
	Update(ctx context.Context, inst *models.Installment) error
	Delete(ctx context.Context, id int64) error
}
