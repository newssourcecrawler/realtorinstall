package models

import "time"

type Commission struct {
	ID               int64     `db:"id" json:"id"`
	TenantID         string    `db:"tenant_id" json:"tenantID"`
	TransactionType  string    `db:"transaction_type" json:"transactiontype"`
	TransactionID    int64     `db:"transaction_id" json:"transactionID"`
	BeneficiaryID    int64     `db:"beneficiary_id" json:"beneficiaryID"`
	CommissionType   string    `db:"commission_type" json:"commissiontype"`
	RateOrAmount     float64   `db:"rate_or_amount" json:"rate_or_amount"`
	CalculatedAmount float64   `db:"calculated_amount" json:"calculatedamount"`
	Memo             string    `db:"memo" json:"memo"`
	CreatedBy        string    `db:"created_by" json:"created_by"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	ModifiedBy       string    `db:"modified_by" json:"modified_by"`
	LastModified     time.Time `db:"last_modified" json:"last_modified"`
	Deleted          bool      `db:"deleted" json:"deleted"`
}

type CommissionSummary struct {
	BeneficiaryID   int64   `json:"beneficiary_id"`
	TotalCommission float64 `json:"total_commission"`
}
