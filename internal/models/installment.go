package models

import "time"

// Installment represents a single payment installment in a plan.
type Installment struct {
	ID             int64     `json:"id"`              // Primary key
	PlanID         int64     `json:"plan_id"`         // FK → InstallmentPlan.ID
	SequenceNumber int       `json:"sequence_number"` // 1, 2, 3, … up to NumInstallments
	DueDate        time.Time `json:"due_date"`        // When this installment is due
	AmountDue      float64   `json:"amount_due"`      // Base amount due (Principal only)
	AmountPaid     float64   `json:"amount_paid"`     // How much has been paid so far
	Status         string    `json:"status"`          // "Pending", "Paid", "Overdue"
	LateFee        float64   `json:"late_fee"`        // Late fee accrued (if overdue)
	PaidDate       time.Time `json:"paid_date"`       // When this installment was fully paid (zero if not paid)
	CreatedAt      time.Time `json:"created_at"`      // When this record was created
	LastModified   time.Time `json:"last_modified"`   // When last updated
}
