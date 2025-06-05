package repos

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/newssourcecrawler/realtorinstall/api/models"
)

type sqliteCommissionRepo struct {
	db *sql.DB
}

func NewSQLiteCommissionRepo(dbPath string) (CommissionRepo, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	schema := `
	CREATE TABLE IF NOT EXISTS commissions (
	  id                INTEGER PRIMARY KEY AUTOINCREMENT,
	  tenant_id         TEXT    NOT NULL,
	  transaction_type  TEXT    NOT NULL,
	  transaction_id    INTEGER NOT NULL,
	  beneficiary_id    INTEGER NOT NULL,
	  commission_type   TEXT    NOT NULL,
	  rate_or_amount    REAL    NOT NULL,
	  calculated_amount REAL    NOT NULL,
	  memo              TEXT,
	  created_by        TEXT    NOT NULL,
	  created_at        DATETIME NOT NULL,
	  modified_by       TEXT    NOT NULL,
	  last_modified     DATETIME NOT NULL,
	  deleted           INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_commissions_tenant ON commissions(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_commissions_txn     ON commissions(transaction_type, transaction_id);
	CREATE INDEX IF NOT EXISTS idx_commissions_benef    ON commissions(beneficiary_id);
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}
	return &sqliteCommissionRepo{db: db}, nil
}

func (r *sqliteCommissionRepo) Create(ctx context.Context, comm *models.Commission) (int64, error) {
	if comm.TenantID == "" ||
		comm.TransactionType == "" ||
		comm.TransactionID == 0 ||
		comm.BeneficiaryID == 0 ||
		comm.CommissionType == "" ||
		comm.CreatedBy == "" ||
		comm.ModifiedBy == "" {
		return 0, errors.New("missing required fields or tenant/audit info")
	}
	now := time.Now().UTC()
	comm.CreatedAt = now
	comm.LastModified = now

	query := `
	INSERT INTO commissions (
	  tenant_id, transaction_type, transaction_id, beneficiary_id,
	  commission_type, rate_or_amount, calculated_amount, memo,
	  created_by, created_at, modified_by, last_modified, deleted
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 0);
	`
	res, err := r.db.ExecContext(ctx, query,
		comm.TenantID,
		comm.TransactionType,
		comm.TransactionID,
		comm.BeneficiaryID,
		comm.CommissionType,
		comm.RateOrAmount,
		comm.CalculatedAmount,
		comm.Memo,
		comm.CreatedBy,
		comm.CreatedAt,
		comm.ModifiedBy,
		comm.LastModified,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *sqliteCommissionRepo) GetByID(ctx context.Context, tenantID string, id int64) (*models.Commission, error) {
	query := `
	SELECT id, tenant_id, transaction_type, transaction_id, beneficiary_id,
	       commission_type, rate_or_amount, calculated_amount, memo,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM commissions
	WHERE tenant_id = ? AND id = ? AND deleted = 0;
	`
	row := r.db.QueryRowContext(ctx, query, tenantID, id)

	var comm models.Commission
	var deletedInt int
	err := row.Scan(
		&comm.ID,
		&comm.TenantID,
		&comm.TransactionType,
		&comm.TransactionID,
		&comm.BeneficiaryID,
		&comm.CommissionType,
		&comm.RateOrAmount,
		&comm.CalculatedAmount,
		&comm.Memo,
		&comm.CreatedBy,
		&comm.CreatedAt,
		&comm.ModifiedBy,
		&comm.LastModified,
		&deletedInt,
	)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	comm.Deleted = (deletedInt != 0)
	return &comm, nil
}

func (r *sqliteCommissionRepo) ListAll(ctx context.Context, tenantID string) ([]*models.Commission, error) {
	query := `
	SELECT id, tenant_id, transaction_type, transaction_id, beneficiary_id,
	       commission_type, rate_or_amount, calculated_amount, memo,
	       created_by, created_at, modified_by, last_modified, deleted
	FROM commissions
	WHERE tenant_id = ? AND deleted = 0;
	`
	rows, err := r.db.QueryContext(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Commission
	for rows.Next() {
		var comm models.Commission
		var deletedInt int
		if err := rows.Scan(
			&comm.ID,
			&comm.TenantID,
			&comm.TransactionType,
			&comm.TransactionID,
			&comm.BeneficiaryID,
			&comm.CommissionType,
			&comm.RateOrAmount,
			&comm.CalculatedAmount,
			&comm.Memo,
			&comm.CreatedBy,
			&comm.CreatedAt,
			&comm.ModifiedBy,
			&comm.LastModified,
			&deletedInt,
		); err != nil {
			return nil, err
		}
		comm.Deleted = (deletedInt != 0)
		out = append(out, &comm)
	}
	return out, nil
}

func (r *sqliteCommissionRepo) Update(ctx context.Context, comm *models.Commission) error {
	existing, err := r.GetByID(ctx, comm.TenantID, comm.ID)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	now := time.Now().UTC()
	comm.LastModified = now

	query := `
	UPDATE commissions
	SET transaction_type = ?, transaction_id = ?, beneficiary_id = ?,
	    commission_type = ?, rate_or_amount = ?, calculated_amount = ?, memo = ?,
	    modified_by = ?, last_modified = ?, deleted = ?
	WHERE tenant_id = ? AND id = ?;
	`
	_, err = r.db.ExecContext(ctx, query,
		comm.TransactionType,
		comm.TransactionID,
		comm.BeneficiaryID,
		comm.CommissionType,
		comm.RateOrAmount,
		comm.CalculatedAmount,
		comm.Memo,
		comm.ModifiedBy,
		comm.LastModified,
		boolToInt(comm.Deleted),
		comm.TenantID,
		comm.ID,
	)
	return err
}

func (r *sqliteCommissionRepo) Delete(ctx context.Context, tenantID string, id int64) error {
	existing, err := r.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if existing.Deleted {
		return ErrNotFound
	}
	query := `
	UPDATE commissions
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
