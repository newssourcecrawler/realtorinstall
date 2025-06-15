package repos

import (
	"context"
	"database/sql"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

// RoleRepo defines role CRUD.
type RoleRepo interface {
	Create(ctx context.Context, r *models.Role) (int64, error)
	ListAll(ctx context.Context) ([]*models.Role, error)
	//GetByID(ctx context.Context, tenantID string, id int64) (*models.Role, error)
	Update(ctx context.Context, b *models.Role) error // using b.TenantID,b.ID
	Delete(ctx context.Context, tenantID string, id int64) error
}

func NewDBRoleRepo(db *sql.DB, driver string) RoleRepo {
	switch driver {
	case "sqlite":
		return NewSQLiteRoleRepo(db)
	case "postgres":
		return NewPostgresRoleRepo(db)
	default:
		panic("unsupported driver for roles: " + driver)
	}
}
