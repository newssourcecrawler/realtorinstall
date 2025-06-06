package models

import "time"

// InstallmentPlan represents a payment plan for a property sale.
type InstallmentPlan struct {
	ID               int64     `db:"id" json:"id"`
	TenantID         string    `db:"tenant_id" json:"tenantID"`
	PropertyID       int64     `db:"property_id" json:"property_id"` // FK → Property.ID
	BuyerID          int64     `db:"buyer_id" json:"buyer_id"`       // FK → Buyer.ID
	TotalPrice       float64   `db:"total_price" json:"total_price"`
	DownPayment      float64   `db:"down_payment" json:"down_payment"`
	NumInstallments  int       `db:"num_installments" json:"num_installments"`
	Frequency        string    `db:"frequency" json:"frequency"` // e.g. "Monthly"
	FirstInstallment time.Time `db:"first_installment" json:"first_installment"`
	InterestRate     float64   `db:"interest_rate" json:"interest_rate"`
	CreatedBy        string    `db:"created_by" json:"created_by"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	ModifiedBy       string    `db:"modified_by" json:"modified_by"`
	LastModified     time.Time `db:"last_modified" json:"last_modified"`
	Deleted          bool      `db:"deleted" json:"deleted"`
}

type PlanSummary struct {
	PlanID           int64   `json:"plan_id"`
	TotalOutstanding float64 `json:"total_outstanding"`
}
