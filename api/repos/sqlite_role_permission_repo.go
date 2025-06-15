package repos

import (
	"context"
	"database/sql"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type sqliteRolePermissionRepo struct{ db *sql.DB }

func NewSQliteRolePermissionRepo(db *sql.DB) RolePermissionRepo {
	return &sqliteRolePermissionRepo{db: db}
}

func (r *sqliteRolePermissionRepo) Create(ctx context.Context, m *models.RolePermission) (int64, error) {
	query := `INSERT INTO roles (name, description) VALUES (?,?)`
	res, err := r.db.ExecContext(ctx, query, m.Name, m.Description)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *sqliteRolePermissionRepo) ListAll(ctx context.Context) ([]*models.RolePermission, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id,name,description FROM roles`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []*models.RolePermission
	for rows.Next() {
		var m models.RolePermission
		rows.Scan(&m.ID, &m.Name, &m.Description)
		out = append(out, &m)
	}
	return out, nil
}

func (r *sqliteRolePermissionRepo) Update(ctx context.Context, m *models.RolePermission) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE roles SET name=?,description=? WHERE id=?`,
		m.Name, m.Description, m.ID,
	)
	return err
}

func (r *sqliteRolePermissionRepo) Delete(ctx context.Context, tenantID string, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM roles WHERE id=?`, id)
	return err
}
