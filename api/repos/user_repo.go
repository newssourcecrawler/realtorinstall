package repos

import (
	"context"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

// UserRepo defines CRUD for User accounts.
type UserRepo interface {
	// Create a new User. Returns the newly‐assigned ID.
	Create(ctx context.Context, u *models.User) (int64, error)

	// GetByID returns a single User by its ID, within the given tenant.
	// Returns ErrNotFound if no such (tenantID, id) record exists or is marked deleted.
	GetByID(ctx context.Context, tenantID string, id int64) (*models.User, error)

	// GetByUsername returns a single User by its username (unique per tenant).
	// Returns ErrNotFound if no such (tenantID, username) record exists or is marked deleted.
	GetByUsername(ctx context.Context, tenantID, username string) (*models.User, error)

	// ListAll returns all non‐deleted Users for a given tenant.
	ListAll(ctx context.Context, tenantID string) ([]*models.User, error)

	// Update modifies an existing User (must have u.TenantID and u.ID set).
	// Returns ErrNotFound if either the row doesn’t exist or is already deleted.
	Update(ctx context.Context, u *models.User) error

	// Delete performs a “soft delete” for (tenantID, id). If that row is already deleted,
	// returns ErrNotFound. Otherwise sets deleted=1, modified_by, last_modified.
	Delete(ctx context.Context, tenantID string, id int64) error
}
