package repos

import (
	"context"
	"database/sql"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type SalesRepo interface {
	Create(ctx context.Context, p *models.Sales) (int64, error) // p.TenantID set
	GetByID(ctx context.Context, tenantID string, id int64) (*models.Sales, error)
	ListAll(ctx context.Context, tenantID string) ([]*models.Sales, error)
	Update(ctx context.Context, b *models.Sales) error // using b.TenantID,b.ID
	Delete(ctx context.Context, tenantID string, id int64) error
	SummarizeByMonth(ctx context.Context, tenantID string) ([]models.MonthSales, error)
}

// MonthSales holds “YYYY‐MM” as Month, plus total sold amount.
type MonthSales struct {
	Month      string  `json:"month"`
	TotalSales float64 `json:"total_sales"`
}

// NewDBSalesRepo selects the concrete implementation based on driver.
func NewDBSalesRepo(db *sql.DB, driver string) SalesRepo {
	switch driver {
	case "postgres":
		return &postgresSalesRepo{db: db}
	case "sqlite":
		return &sqliteSalesRepo{db: db}
	default:
		panic("unsupported driver: " + driver)
	}
}
