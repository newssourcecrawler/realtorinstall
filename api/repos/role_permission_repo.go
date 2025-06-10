package repos

import (
	"context"
	"database/sql"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

// RolePermissionRepo defines role–permission assignments.
type RolePermissionRepo interface {
	Add(ctx context.Context, rp *models.RolePermission) error
	ListByRole(ctx context.Context, roleID int64) ([]int64 /*permissionIDs*/, error)
	Create(ctx context.Context, r *models.RolePermission) (int64, error)
	ListAll(ctx context.Context) ([]*models.RolePermission, error)
	GetByID(ctx context.Context, tenantID string, id int64) (*models.RolePermission, error)
	Update(ctx context.Context, b *models.RolePermission) error // using b.TenantID,b.ID
	Delete(ctx context.Context, tenantID string, id int64) error
}

func NewDBRolePermissionRepo(db *sql.DB, driver string) RolePermissionRepo {
	switch driver {
	case "sqlite":
		return NewSQLiteRolePermissionRepo(db)
	case "postgres":
		return NewPostgresRolePermissionRepo(db)
	default:
		panic("unsupported driver for roles: " + driver)
	}
}
