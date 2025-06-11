package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type sqliteSalesRepo struct {
	db *sql.DB
}

func NewSQLiteSalesRepo(dbPath string) (SalesRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	schema := `
	CREATE TABLE IF NOT EXISTS sales (
	  id            INTEGER PRIMARY KEY AUTOINCREMENT,
	  tenant_id     TEXT    NOT NULL,
	  property_id   INTEGER NOT NULL,
	  buyer_id      INTEGER NOT NULL,
	  sale_price    REAL    NOT NULL,
	  sale_date     DATETIME NOT NULL,
	  sale_type     TEXT    NOT NULL,
	  created_by    TEXT    NOT NULL,
	  created_at    DATETIME NOT NULL,
	  modified_by   TEXT    NOT NULL,
	  last_modified DATETIME NOT NULL,
	  deleted       INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_sales_tenant     ON sales(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_sales_property   ON sales(property_id);
	CREATE INDEX IF NOT EXISTS idx_sales_buyer      ON sales(buyer_id);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &sqliteSalesRepo{db: db}, nil
}

func (r *sqliteSalesRepo) Create(ctx context.Context, s *models.Sales) (int64, error) {
	if s.TenantID == "" || s.PropertyID == 0 || s.BuyerID == 0 ||
		s.SalePrice <= 0 || s.SaleType == "" ||
		s.CreatedBy == "" || s.ModifiedBy == "" {
		return 0, errors.New("missing required fields or tenant/audit info")
	}
	now := time.Now().UTC()
	s.CreatedAt = now
	s.LastModified = now

	query := `
	INSERT INTO sales (
	  tenant_id, property_id, buyer_id, sale_price, sale_date, sale_type,
	  created_by, created_at, modified_by, last_modified, deleted
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0);
	`
	res, err := r.db.ExecContext(ctx, query,
		s.TenantID,
		s.PropertyID,
		s.BuyerID,
		s.SalePrice,
		s.SaleDate,
		s.SaleType,
		s.CreatedBy,
		s.CreatedAt,
		s.ModifiedBy,
		s.LastModified,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *sqliteSalesRepo) GetByID(ctx context.Context, tenantID string, id int64) (*models.Sales, error) {
	query := `
	SELECT id, tenant_id, property_id, buyer_id, sale_price, sale_date, sale_type,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM sales
	WHERE tenant_id = ? AND id = ? AND deleted = 0;
	`
	row := r.db.QueryRowContext(ctx, query, tenantID, id)

	var s models.Sales
	var deletedInt int
	err := row.Scan(
		&s.ID,
		&s.TenantID,
		&s.PropertyID,
		&s.BuyerID,
		&s.SalePrice,
		&s.SaleDate,
		&s.SaleType,
		&s.CreatedBy,
		&s.CreatedAt,
		&s.ModifiedBy,
		&s.LastModified,
		&deletedInt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	s.Deleted = (deletedInt != 0)
	return &s, nil
}

func (r *sqliteSalesRepo) ListAll(ctx context.Context, tenantID string) ([]*models.Sales, error) {
	query := `
	SELECT id, tenant_id, property_id, buyer_id, sale_price, sale_date, sale_type,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM sales
	WHERE tenant_id = ? AND deleted = 0;
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Sales
	for rows.Next() {
		var s models.Sales
		var deletedInt int
		if err := rows.Scan(
			&s.ID,
			&s.TenantID,
			&s.PropertyID,
			&s.BuyerID,
			&s.SalePrice,
			&s.SaleDate,
			&s.SaleType,
			&s.CreatedBy,
			&s.CreatedAt,
			&s.ModifiedBy,
			&s.LastModified,
			&deletedInt,
		); err != nil {
			return nil, err
		}
		s.Deleted = (deletedInt != 0)
		out = append(out, &s)
	}
	return out, nil
}

func (r *sqliteSalesRepo) Update(ctx context.Context, s *models.Sales) error {
	existing, err := r.GetByID(ctx, s.TenantID, s.ID)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	now := time.Now().UTC()
	s.LastModified = now

	query := `
	UPDATE sales
	SET property_id = ?, buyer_id = ?, sale_price = ?, sale_date = ?, sale_type = ?,
	    modified_by = ?, last_modified = ?, deleted = ?
	WHERE tenant_id = ? AND id = ?;
	`
	_, err = r.db.ExecContext(ctx, query,
		s.PropertyID,
		s.BuyerID,
		s.SalePrice,
		s.SaleDate,
		s.SaleType,
		s.ModifiedBy,
		s.LastModified,
		boolToInt(s.Deleted),
		s.TenantID,
		s.ID,
	)
	return err
}

func (r *sqliteSalesRepo) Delete(ctx context.Context, tenantID string, id int64) error {
	existing, err := r.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	query := `
	UPDATE sales
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

func (r *sqliteCommissionRepo) SQLiteSummarizeByBeneficiary(ctx context.Context, tenantID string) ([]models.CommissionSummary, error) {
	query := `
        SELECT beneficiary_id, SUM(calculated_amount) AS total_commission
          FROM commissions
         WHERE tenant_id = ? AND deleted = 0
         GROUP BY beneficiary_id;
    `
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.CommissionSummary
	for rows.Next() {
		var cs models.CommissionSummary
		if err := rows.Scan(&cs.BeneficiaryID, &cs.TotalCommission); err != nil {
			return nil, err
		}
		out = append(out, cs)
	}
	return out, nil
}

// SummarizeByMonth returns sum of sale_price per month (YYYY‐MM).
func (r *sqliteSalesRepo) SummarizeByMonth(ctx context.Context, tenantID string) ([]models.MonthSales, error) {
	query := `
        SELECT strftime('%Y-%m', sale_date) AS month,
               SUM(sale_price)       AS total_sales
          FROM sales
         WHERE tenant_id = ? AND deleted = 0
         GROUP BY month
         ORDER BY month DESC;
    `
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.MonthSales
	for rows.Next() {
		var ms models.MonthSales
		if err := rows.Scan(&ms.Month, &ms.TotalSales); err != nil {
			return nil, err
		}
		out = append(out, ms)
	}
	return out, nil
}
