package models

import "time"

// Payment represents a payment made against an installment.
type Payment struct {
	ID             int64     `db:"id" json:"id"`
	TenantID       string    `db:"tenant_id" json:"tenantID"`
	InstallmentID  int64     `db:"installment_id" json:"installment_id"` // FK → Installment.ID
	AmountPaid     float64   `db:"amount_paid" json:"amount_paid"`
	PaymentDate    time.Time `db:"payment_date" json:"payment_date"`
	PaymentMethod  string    `db:"payment_method" json:"payment_method"`
	TransactionRef string    `db:"transaction_ref" json:"transaction_ref"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	ModifiedBy     string    `db:"modified_by" json:"modified_by"`
	LastModified   time.Time `db:"last_modified" json:"last_modified"`
	Deleted        bool      `db:"deleted" json:"deleted"`
}
