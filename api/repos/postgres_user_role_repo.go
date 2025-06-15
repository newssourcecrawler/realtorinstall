package repos

import (
	"context"
	"database/sql"

	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type postgresUserRoleRepo struct{ db *sql.DB }

func NewPostgresUserRoleRepo(db *sql.DB) RoleRepo {
	return &postgresUserRoleRepo{db: db}
}

func (r *postgresUserRoleRepo) Create(ctx context.Context, m *models.Role) (int64, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO roles (name, description) VALUES ($1,$2) RETURNING id`,
		m.Name, m.Description,
	).Scan(&id)
	return id, err
}

func (r *postgresUserRoleRepo) ListAll(ctx context.Context) ([]*models.Role, error) {
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

func (r *postgresUserRoleRepo) Update(ctx context.Context, m *models.Role) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE roles SET name=$1,description=$2 WHERE id=$3`,
		m.Name, m.Description, m.ID,
	)
	return err
}

func (r *postgresUserRoleRepo) Delete(ctx context.Context, tenantID string, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM roles WHERE id=$1`, id)
	return err
}
