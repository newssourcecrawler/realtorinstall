package models

import "time"

// Installment represents a single payment installment in a plan.
type Installment struct {
	ID             int64     `db:"id" json:"id"`
	TenantID       string    `db:"tenant_id" json:"tenantID"`
	PlanID         int64     `db:"plan_id" json:"plan_id"`                 // FK → InstallmentPlan.ID
	SequenceNumber int       `db:"sequence_number" json:"sequence_number"` // 1…NumInstallments
	DueDate        time.Time `db:"due_date" json:"due_date"`
	AmountDue      float64   `db:"amount_due" json:"amount_due"`
	AmountPaid     float64   `db:"amount_paid" json:"amount_paid"`
	Status         string    `db:"status" json:"status"` // "Pending", "Paid", "Overdue"
	LateFee        float64   `db:"late_fee" json:"late_fee"`
	PaidDate       time.Time `db:"paid_date" json:"paid_date"`
	CreatedBy      string    `db:"created_by" json:"created_by"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	ModifiedBy     string    `db:"modified_by" json:"modified_by"`
	LastModified   time.Time `db:"last_modified" json:"last_modified"`
	Deleted        bool      `db:"deleted" json:"deleted"`
}
