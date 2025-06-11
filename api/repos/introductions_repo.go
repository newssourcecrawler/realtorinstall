package repos

import (
	"context"
	"database/sql"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type IntroductionsRepo interface {
	Create(ctx context.Context, p *models.Introductions) (int64, error) // p.TenantID set
	GetByID(ctx context.Context, tenantID string, id int64) (*models.Introductions, error)
	ListAll(ctx context.Context, tenantID string) ([]*models.Introductions, error)
	Update(ctx context.Context, b *models.Introductions) error // using b.TenantID,b.ID
	Delete(ctx context.Context, tenantID string, id int64) error
}

// NewDBIntroductionsRepo selects the concrete implementation based on driver.
func NewDBIntroductionsRepo(db *sql.DB, driver string) IntroductionsRepo {
	switch driver {
	case "postgres":
		return &postgresIntroductionsRepo{db: db}
	case "sqlite":
		return &sqliteIntroductionsRepo{db: db}
	default:
		panic("unsupported driver: " + driver)
	}
}
