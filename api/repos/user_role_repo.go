package repos

import (
	"context"
	"database/sql"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

// UserRoleRepo defines user–role assignments.
type UserRoleRepo interface {
	Add(ctx context.Context, ur *models.UserRole) error
	ListRoles(ctx context.Context, userID int64) ([]int64 /*roleIDs*/, error)
	Create(ctx context.Context, r *models.UserRole) (int64, error)
	ListAll(ctx context.Context) ([]*models.UserRole, error)
	GetByID(ctx context.Context, tenantID string, id int64) (*models.UserRole, error)
	Update(ctx context.Context, b *models.UserRole) error // using b.TenantID,b.ID
	Delete(ctx context.Context, tenantID string, id int64) error
}

func NewDBUserRoleRepo(db *sql.DB, driver string) UserRoleRepo {
	switch driver {
	case "postgres":
		return NewPostgresUserRoleRepo(db)
	case "sqlite":
		return NewSQLiteUserRoleRepo(db)
	default:
		panic("unsupported driver: " + driver)
	}
}
