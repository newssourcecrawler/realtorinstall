package repos

import (
	"context"
	"database/sql"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

// PermissionRepo defines permission CRUD.
type PermissionRepo interface {
	Create(ctx context.Context, p *models.Permission) (int64, error)
	ListAll(ctx context.Context) ([]*models.Permission, error)
	Update(ctx context.Context, b *models.Permission) error
	Delete(ctx context.Context, tenantID string, id int64) error
}

func NewDBPermissionRepo(db *sql.DB, driver string) PermissionRepo {
	switch driver {
	case "sqlite":
		return NewSQLitePermissionRepo(db)
	case "postgres":
		return NewPostgresPermissionRepo(db)
	default:
		panic("unsupported driver for roles: " + driver)
	}
}
