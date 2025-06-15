package repos

import (
	"context"
	"database/sql"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type sqliteRoleRepo struct{ db *sql.DB }

func NewSQLiteRoleRepo(db *sql.DB) RoleRepo {
	return &sqliteRoleRepo{db: db}
}

func (r *sqliteRoleRepo) Create(ctx context.Context, m *models.Role) (int64, error) {
	query := `INSERT INTO roles (name, description) VALUES (?,?)`
	res, err := r.db.ExecContext(ctx, query, m.Name, m.Description)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *sqliteRoleRepo) ListAll(ctx context.Context) ([]*models.Role, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id,name,description FROM roles`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.Role
	for rows.Next() {
		var m models.Role
		rows.Scan(&m.ID, &m.Name, &m.Description)
		out = append(out, &m)
	}
	return out, nil
}

func (r *sqliteRoleRepo) Update(ctx context.Context, m *models.Role) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE roles SET name=?,description=? WHERE id=?`,
		m.Name, m.Description, m.ID,
	)
	return err
}

func (r *sqliteRoleRepo) Delete(ctx context.Context, tenantID string, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM roles WHERE id=?`, id)
	return err
}
