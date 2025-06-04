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
	// Create users table (with soft-delete)
	schema := `
    CREATE TABLE IF NOT EXISTS users (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      username TEXT NOT NULL UNIQUE,
      password_hash TEXT NOT NULL,
	  first_name TEXT NOT NULL,
	  last_name TEXT NOT NULL,
      role TEXT NOT NULL,
	  email TEXT NOT NULL,
	  phone TEXT NOT NULL,
      created_at DATETIME NOT NULL,
	  created_by TEXT NOT NULL,
      last_modified DATETIME NOT NULL,
	  modified_by TEXT NOT NULL,
      deleted INTEGER NOT NULL DEFAULT 0
    );
    `
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &sqliteUserRepo{db: db}, nil
}

func (r *sqliteUserRepo) Create(ctx context.Context, u *models.User) (int64, error) {
	if u.UserName == "" || u.PasswordHash == "" || u.Role == "" || u.FirstName == "" || u.LastName == "" {
		return 0, errors.New("username, password, first and last names and role required")
	}
	now := time.Now().UTC()
	u.CreatedAt = now
	u.LastModified = now

	query := `
    INSERT INTO users (
      username, password_hash, role, first_name, last_name, created_at, last_modified, deleted
    ) VALUES (?, ?, ?, ?, ?, 0);
    `
	res, err := r.db.ExecContext(ctx, query,
		u.UserName,
		u.PasswordHash,
		u.Role,
		u.FirstName,
		u.LastName,
		u.CreatedAt,
		u.LastModified,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *sqliteUserRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
    SELECT id, username, password_hash, role, first_name, last_name, created_at, last_modified, deleted
    FROM users
    WHERE id = ?;
    `
	row := r.db.QueryRowContext(ctx, query, id)
	var u models.User
	var deletedInt int
	err := row.Scan(
		&u.ID,
		u.UserName,
		u.PasswordHash,
		u.Role,
		u.FirstName,
		u.LastName,
		u.CreatedAt,
		u.LastModified,
		&deletedInt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *sqliteUserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
    SELECT id, username, password_hash, role, first_name, last_name, created_at, last_modified, deleted
    FROM users
    WHERE username = ? AND deleted = 0;
    `
	row := r.db.QueryRowContext(ctx, query, username)
	var u models.User
	var deletedInt int
	err := row.Scan(
		&u.ID,
		u.UserName,
		u.PasswordHash,
		u.Role,
		u.FirstName,
		u.LastName,
		u.CreatedAt,
		u.LastModified,
		&deletedInt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// ListAll, Update, Delete follow same pattern (soft-deleting by flagging deleted=1)
