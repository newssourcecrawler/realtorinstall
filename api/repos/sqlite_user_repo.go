package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type sqliteUserRepo struct {
	db *sql.DB
}

func NewSQLiteUserRepo(dbPath string) (UserRepo, error) {
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
	return &sqliteUserRepo{db: db}, nil
}

func (r *sqliteUserRepo) Create(ctx context.Context, u *models.User) (int64, error) {
	if u.TenantID == "" ||
		u.UserName == "" ||
		u.PasswordHash == "" ||
		u.FirstName == "" ||
		u.LastName == "" ||
		u.Role == "" ||
		u.Email == "" ||
		u.CreatedBy == "" ||
		u.ModifiedBy == "" {
		return 0, errors.New("missing required fields or tenant/audit info")
	}
	now := time.Now().UTC()
	u.CreatedAt = now
	u.LastModified = now

	query := `
	INSERT INTO users (
	  tenant_id, username, password_hash, first_name, last_name, role, email, phone,
	  created_by, created_at, modified_by, last_modified, deleted
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0);
	`
	res, err := r.db.ExecContext(ctx, query,
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
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *sqliteUserRepo) GetByID(ctx context.Context, tenantID string, id int64) (*models.User, error) {
	query := `
	SELECT id, tenant_id, username, password_hash, first_name, last_name, role, email, phone,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM users
	WHERE tenant_id = ? AND id = ? AND deleted = 0;
	`
	row := r.db.QueryRowContext(ctx, query, tenantID, id)

	var u models.User
	var deletedInt int
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
		&deletedInt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	u.Deleted = (deletedInt != 0)
	return &u, nil
}

func (r *sqliteUserRepo) GetByUsername(ctx context.Context, tenantID, username string) (*models.User, error) {
	query := `
	SELECT id, tenant_id, username, password_hash, first_name, last_name, role, email, phone,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM users
	WHERE tenant_id = ? AND username = ? AND deleted = 0;
	`
	row := r.db.QueryRowContext(ctx, query, tenantID, username)

	var u models.User
	var deletedInt int
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
		&deletedInt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	u.Deleted = (deletedInt != 0)
	return &u, nil
}

func (r *sqliteUserRepo) ListAll(ctx context.Context, tenantID string) ([]*models.User, error) {
	query := `
	SELECT id, tenant_id, username, password_hash, first_name, last_name, role, email, phone,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM users
	WHERE tenant_id = ? AND deleted = 0;
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.User
	for rows.Next() {
		var u models.User
		var deletedInt int
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
			&deletedInt,
		); err != nil {
			return nil, err
		}
		u.Deleted = (deletedInt != 0)
		out = append(out, &u)
	}
	return out, nil
}

func (r *sqliteUserRepo) Update(ctx context.Context, u *models.User) error {
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
	UPDATE users
	SET username = ?, password_hash = ?, first_name = ?, last_name = ?, role = ?, email = ?, phone = ?, modified_by = ?, last_modified = ?, deleted = ?
	WHERE tenant_id = ? AND id = ?;
	`
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
		boolToInt(u.Deleted),
		u.TenantID,
		u.ID,
	)
	return err
}

func (r *sqliteUserRepo) Delete(ctx context.Context, tenantID string, id int64) error {
	existing, err := r.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	query := `
	UPDATE users
	SET deleted = 1, modified_by = ?, last_modified = ?
	WHERE tenant_id = ? AND id = ?;
	`
	_, err = r.db.ExecContext(ctx, query,
		existing.ModifiedBy,
		time.Now().UTC(),
		tenantID,
		id,
	)
	return err
}
