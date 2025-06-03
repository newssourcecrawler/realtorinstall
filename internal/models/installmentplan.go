package models

import "time"

// InstallmentPlan represents a payment plan for a property sale.
type InstallmentPlan struct {
	ID               int64     `json:"id"`                // Primary key
	PropertyID       int64     `json:"property_id"`       // FK → Property.ID
	BuyerID          int64     `json:"buyer_id"`          // FK → Buyer.ID
	TotalPrice       float64   `json:"total_price"`       // Total sale price
	DownPayment      float64   `json:"down_payment"`      // Amount due at booking
	NumInstallments  int       `json:"num_installments"`  // Number of equal installments
	Frequency        string    `json:"frequency"`         // "Monthly", "Quarterly", etc.
	FirstInstallment time.Time `json:"first_installment"` // Date of installment #1
	InterestRate     float64   `json:"interest_rate"`     // Annual interest % (if any)
	CreatedAt        time.Time `json:"created_at"`        // When plan was created
	LastModified     time.Time `json:"last_modified"`     // When plan was last modified
}
