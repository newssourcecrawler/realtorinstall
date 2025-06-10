package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type sqliteLettingsRepo struct {
	db *sql.DB
}

func NewSQLiteLettingsRepo(dbPath string) (LettingsRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	schema := `
	CREATE TABLE IF NOT EXISTS lettings (
	  id             INTEGER PRIMARY KEY AUTOINCREMENT,
	  tenant_id      TEXT    NOT NULL,
	  property_id    INTEGER NOT NULL,
	  tenant_user_id INTEGER NOT NULL,
	  rent_amount    REAL    NOT NULL,
	  rent_term      INTEGER NOT NULL,
	  rent_cycle     TEXT    NOT NULL,
	  memo           TEXT,
	  start_date     DATETIME NOT NULL,
	  end_date       DATETIME,
	  created_by     TEXT    NOT NULL,
	  created_at     DATETIME NOT NULL,
	  modified_by    TEXT    NOT NULL,
	  last_modified  DATETIME NOT NULL,
	  deleted        INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_lettings_tenant    ON lettings(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_lettings_property  ON lettings(property_id);
	CREATE INDEX IF NOT EXISTS idx_lettings_tenantuser ON lettings(tenant_user_id);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &sqliteLettingsRepo{db: db}, nil
}

func (r *sqliteLettingsRepo) Create(ctx context.Context, lt *models.Lettings) (int64, error) {
	if lt.TenantID == "" ||
		lt.PropertyID == 0 ||
		lt.TenantUserID == 0 ||
		lt.RentAmount <= 0 ||
		lt.RentTerm <= 0 ||
		lt.RentCycle == "" ||
		lt.StartDate.IsZero() ||
		lt.CreatedBy == "" ||
		lt.ModifiedBy == "" {
		return 0, errors.New("missing required fields or tenant/audit info")
	}
	now := time.Now().UTC()
	lt.CreatedAt = now
	lt.LastModified = now

	query := `
	INSERT INTO lettings (
	  tenant_id, property_id, tenant_user_id, rent_amount, rent_term, rent_cycle, memo, start_date, end_date,
	  created_by, created_at, modified_by, last_modified, deleted
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0);
	`
	res, err := r.db.ExecContext(ctx, query,
		lt.TenantID,
		lt.PropertyID,
		lt.TenantUserID,
		lt.RentAmount,
		lt.RentTerm,
		lt.RentCycle,
		lt.Memo,
		lt.StartDate,
		lt.EndDate,
		lt.CreatedBy,
		lt.CreatedAt,
		lt.ModifiedBy,
		lt.LastModified,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *sqliteLettingsRepo) GetByID(ctx context.Context, tenantID string, id int64) (*models.Lettings, error) {
	query := `
	SELECT id, tenant_id, property_id, tenant_user_id, rent_amount, rent_term, rent_cycle, memo, start_date, end_date,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM lettings
	WHERE tenant_id = ? AND id = ? AND deleted = 0;
	`
	row := r.db.QueryRowContext(ctx, query, tenantID, id)

	var lt models.Lettings
	var deletedInt int
	err := row.Scan(
		&lt.ID,
		&lt.TenantID,
		&lt.PropertyID,
		&lt.TenantUserID,
		&lt.RentAmount,
		&lt.RentTerm,
		&lt.RentCycle,
		&lt.Memo,
		&lt.StartDate,
		&lt.EndDate,
		&lt.CreatedBy,
		&lt.CreatedAt,
		&lt.ModifiedBy,
		&lt.LastModified,
		&deletedInt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	lt.Deleted = (deletedInt != 0)
	return &lt, nil
}

func (r *sqliteLettingsRepo) ListAll(ctx context.Context, tenantID string) ([]*models.Lettings, error) {
	query := `
	SELECT id, tenant_id, property_id, tenant_user_id, rent_amount, rent_term, rent_cycle, memo, start_date, end_date,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM lettings
	WHERE tenant_id = ? AND deleted = 0;
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Lettings
	for rows.Next() {
		var lt models.Lettings
		var deletedInt int
		if err := rows.Scan(
			&lt.ID,
			&lt.TenantID,
			&lt.PropertyID,
			&lt.TenantUserID,
			&lt.RentAmount,
			&lt.RentTerm,
			&lt.RentCycle,
			&lt.Memo,
			&lt.StartDate,
			&lt.EndDate,
			&lt.CreatedBy,
			&lt.CreatedAt,
			&lt.ModifiedBy,
			&lt.LastModified,
			&deletedInt,
		); err != nil {
			return nil, err
		}
		lt.Deleted = (deletedInt != 0)
		out = append(out, &lt)
	}
	return out, nil
}

func (r *sqliteLettingsRepo) Update(ctx context.Context, lt *models.Lettings) error {
	existing, err := r.GetByID(ctx, lt.TenantID, lt.ID)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	now := time.Now().UTC()
	lt.LastModified = now

	query := `
	UPDATE lettings
	SET property_id = ?, tenant_user_id = ?, rent_amount = ?, rent_term = ?, rent_cycle = ?, memo = ?, start_date = ?, end_date = ?,
	    modified_by = ?, last_modified = ?, deleted = ?
	WHERE tenant_id = ? AND id = ?;
	`
	_, err = r.db.ExecContext(ctx, query,
		lt.PropertyID,
		lt.TenantUserID,
		lt.RentAmount,
		lt.RentTerm,
		lt.RentCycle,
		lt.Memo,
		lt.StartDate,
		lt.EndDate,
		lt.ModifiedBy,
		lt.LastModified,
		boolToInt(lt.Deleted),
		lt.TenantID,
		lt.ID,
	)
	return err
}

func (r *sqliteLettingsRepo) Delete(ctx context.Context, lt *models.Lettings) error {
	existing, err := r.GetByID(ctx, lt.TenantID, lt.ID)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	query := `
	UPDATE lettings
	SET deleted = 1, modified_by = ?, last_modified = ?
	WHERE tenant_id = ? AND id = ?;
	`
	_, err = r.db.ExecContext(ctx, query,
		lt.ModifiedBy,
		time.Now().UTC(),
		lt.TenantUserID,
		lt.ID,
	)
	return err
}

// SummarizeRentRoll sums rent_amount for all “currently active” lettings.
func (r *sqliteLettingsRepo) SummarizeRentRoll(ctx context.Context, tenantID string) ([]models.RentRoll, error) {
	query := `
        SELECT property_id, SUM(rent_amount) AS total_rent
          FROM lettings
         WHERE tenant_id = ? 
           AND deleted = 0
           AND DATE(start_date) <= DATE('now')
           AND (end_date IS NULL OR DATE(end_date) > DATE('now'))
         GROUP BY property_id;
    `
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.RentRoll
	for rows.Next() {
		var rr models.RentRoll
		if err := rows.Scan(&rr.PropertyID, &rr.TotalRent); err != nil {
			return nil, err
		}
		out = append(out, rr)
	}
	return out, nil
}
