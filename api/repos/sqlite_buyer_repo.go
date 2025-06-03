package repos

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/newssourcecrawler/realtorinstall/internal/models"
)

// sqliteBuyerRepo is a stub implementation of BuyerRepo.
// You can flesh out these methods later when you add Buyer logic.
type sqliteBuyerRepo struct {
	db *sql.DB
}

// NewSQLiteBuyerRepo opens (or creates) the SQLite file and ensures the "buyers" table exists.
func NewSQLiteBuyerRepo(dbPath string) (BuyerRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create a minimal "buyers" table; adjust columns later as needed.
	schema := `
	CREATE TABLE IF NOT EXISTS buyers (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  name TEXT NOT NULL,
	  email TEXT,
	  phone TEXT
	);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	return &sqliteBuyerRepo{db: db}, nil
}

// The following methods satisfy the BuyerRepo interface but return "not implemented" for now.

func (r *sqliteBuyerRepo) Create(ctx context.Context, b *models.Buyer) (int64, error) {
	return 0, nil // stub
}

func (r *sqliteBuyerRepo) GetByID(ctx context.Context, id int64) (*models.Buyer, error) {
	return nil, nil // stub
}

func (r *sqliteBuyerRepo) ListAll(ctx context.Context) ([]*models.Buyer, error) {
	return nil, nil // stub
}

func (r *sqliteBuyerRepo) Update(ctx context.Context, b *models.Buyer) error {
	return nil // stub
}

func (r *sqliteBuyerRepo) Delete(ctx context.Context, id int64) error {
	return nil // stub
}
