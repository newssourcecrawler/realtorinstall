package repos

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/newssourcecrawler/realtorinstall/api/models"
)

// postgresUserRepo implements UserRepo for PostgreSQL.
type postgresUserRepo struct {
	db *sql.DB
}

// NewPostgresUserRepo returns a new UserRepo backed by the given *sql.DB.
func NewPostgresUserRepo(db *sql.DB) UserRepo {
	return &postgresUserRepo{db: db}
}

/*func NewPostgresUserRepo(dbPath string) (UserRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	schema := `
	CREATE TABLE IF NOT EXISTS users (
	  id            INTEGER PRIMARY KEY AUTOINCREMENT,
	  tenant_id     TEXT    NOT NULL,
	  username      TEXT    NOT NULL UNIQUE,
	  password_hash TEXT    NOT NULL,
	  first_name    TEXT    NOT NULL,
	  last_name     TEXT    NOT NULL,
	  role          TEXT    NOT NULL,
	  email         TEXT    NOT NULL,
	  phone         TEXT,
	  created_by    TEXT    NOT NULL,
	  created_at    DATETIME NOT NULL,
	  modified_by   TEXT    NOT NULL,
	  last_modified DATETIME NOT NULL,
	  deleted       INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_users_tenant ON users(tenant_id);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &postgresUserRepo{db: db}, nil
}
*/

// Create inserts a new user and returns its generated ID.
func (r *postgresUserRepo) Create(ctx context.Context, u *models.User) (int64, error) {
	if u.TenantID == "" || u.UserName == "" || u.PasswordHash == "" ||
		u.FirstName == "" || u.LastName == "" || u.Role == "" ||
		u.Email == "" || u.CreatedBy == "" || u.ModifiedBy == "" {
		return 0, errors.New("missing required fields or tenant/audit info")
	}

	now := time.Now().UTC()
	u.CreatedAt = now
	u.LastModified = now

	query := `
	INSERT INTO users (
		tenant_id, username, password_hash,
		first_name, last_name, role, email, phone,
		created_by, created_at, modified_by, last_modified, deleted
	) VALUES (
		$1,$2,$3,
		$4,$5,$6,$7,$8,
		$9,$10,$11,$12, FALSE
	)
	RETURNING id` // boolean column 'deleted'

	var newID int64
	err := r.db.QueryRowContext(
		ctx, query,
		u.TenantID,
		u.UserName,
		u.PasswordHash,
		u.FirstName,
		u.LastName,
		u.Role,
		u.Email,
		u.Phone,
		u.CreatedBy,
		u.CreatedAt,
		u.ModifiedBy,
		u.LastModified,
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("Create user: %w", err)
	}
	return newID, nil
}

// GetByID fetches a single user by tenantID and id.
func (r *postgresUserRepo) GetByID(ctx context.Context, tenantID string, id int64) (*models.User, error) {
	query := `
	SELECT id, tenant_id, username, password_hash,
		first_name, last_name, role, email, phone,
		created_by, created_at, modified_by, last_modified, deleted
	FROM users
	WHERE tenant_id = $1 AND id = $2 AND deleted = FALSE` + `;`

	row := r.db.QueryRowContext(ctx, query, tenantID, id)
	var u models.User
	var deleted bool
	err := row.Scan(
		&u.ID,
		&u.TenantID,
		&u.UserName,
		&u.PasswordHash,
		&u.FirstName,
		&u.LastName,
		&u.Role,
		&u.Email,
		&u.Phone,
		&u.CreatedBy,
		&u.CreatedAt,
		&u.ModifiedBy,
		&u.LastModified,
		&deleted,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("GetByID user: %w", err)
	}
	u.Deleted = deleted
	return &u, nil
}

// GetByUsername fetches a user by tenantID and username.
func (r *postgresUserRepo) GetByUsername(ctx context.Context, tenantID, username string) (*models.User, error) {
	query := `
	SELECT id, tenant_id, username, password_hash,
		first_name, last_name, role, email, phone,
		created_by, created_at, modified_by, last_modified, deleted
	FROM users
	WHERE tenant_id = $1 AND username = $2 AND deleted = FALSE` + `;`

	row := r.db.QueryRowContext(ctx, query, tenantID, username)
	var u models.User
	var deleted bool
	err := row.Scan(
		&u.ID,
		&u.TenantID,
		&u.UserName,
		&u.PasswordHash,
		&u.FirstName,
		&u.LastName,
		&u.Role,
		&u.Email,
		&u.Phone,
		&u.CreatedBy,
		&u.CreatedAt,
		&u.ModifiedBy,
		&u.LastModified,
		&deleted,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("GetByUsername user: %w", err)
	}
	u.Deleted = deleted
	return &u, nil
}

// ListAll returns all non-deleted users for a tenant.
func (r *postgresUserRepo) ListAll(ctx context.Context, tenantID string) ([]*models.User, error) {
	query := `
	SELECT id, tenant_id, username, password_hash,
		first_name, last_name, role, email, phone,
		created_by, created_at, modified_by, last_modified, deleted
	FROM users
	WHERE tenant_id = $1 AND deleted = FALSE` + `;`

	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, fmt.Errorf("ListAll users: %w", err)
	}
	defer rows.Close()

	var out []*models.User
	for rows.Next() {
		var u models.User
		var deleted bool
		if err := rows.Scan(
			&u.ID,
			&u.TenantID,
			&u.UserName,
			&u.PasswordHash,
			&u.FirstName,
			&u.LastName,
			&u.Role,
			&u.Email,
			&u.Phone,
			&u.CreatedBy,
			&u.CreatedAt,
			&u.ModifiedBy,
			&u.LastModified,
			&deleted,
		); err != nil {
			return nil, fmt.Errorf("ListAll scan: %w", err)
		}
		u.Deleted = deleted
		out = append(out, &u)
	}
	return out, nil
}

// Update modifies an existing, non-deleted user.
func (r *postgresUserRepo) Update(ctx context.Context, u *models.User) error {
	existing, err := r.GetByID(ctx, u.TenantID, u.ID)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	now := time.Now().UTC()
	u.LastModified = now

	query := `
	UPDATE users SET
		username = $1, password_hash = $2,
		first_name = $3, last_name = $4, role = $5,
		email = $6, phone = $7,
		modified_by = $8, last_modified = $9,
		deleted = $10
	WHERE tenant_id = $11 AND id = $12` + `;`

	_, err = r.db.ExecContext(ctx, query,
		u.UserName,
		u.PasswordHash,
		u.FirstName,
		u.LastName,
		u.Role,
		u.Email,
		u.Phone,
		u.ModifiedBy,
		u.LastModified,
		u.Deleted,
		u.TenantID,
		u.ID,
	)
	return err
}

// Delete marks a user as deleted (soft delete).
func (r *postgresUserRepo) Delete(ctx context.Context, tenantID string, id int64) error {
	existing, err := r.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}

	query := `
	UPDATE users SET deleted = TRUE,
		modified_by = $1,
		last_modified = $2
	WHERE tenant_id = $3 AND id = $4` + `;`

	_, err = r.db.ExecContext(ctx, query,
		existing.ModifiedBy,
		time.Now().UTC(),
		tenantID,
		id,
	)
	return err
}

// ListPermissionsForUser returns all permission names for a user.
func (r *postgresUserRepo) ListPermissionsForUser(ctx context.Context, tenantID string, userID int64) ([]string, error) {
	query := `
	SELECT p.name
	FROM permissions p
	JOIN role_permissions rp ON rp.permission_id = p.id
	JOIN user_roles ur ON ur.role_id = rp.role_id
	JOIN users u ON u.id = ur.user_id
	WHERE u.tenant_id = $1 AND u.id = $2 AND u.deleted = FALSE` + `;`

	rows, err := r.db.QueryContext(ctx, query, tenantID, userID)
	if err != nil {
		return nil, fmt.Errorf("ListPermissions: %w", err)
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("ListPermissions scan: %w", err)
		}
		perms = append(perms, name)
	}
	return perms, nil
}
