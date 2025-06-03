package models

import "time"

// Payment represents a payment made against an installment.
type Payment struct {
	ID             int64     `json:"id"`              // Primary key
	InstallmentID  int64     `json:"installment_id"`  // FK → Installment.ID
	AmountPaid     float64   `json:"amount_paid"`     // How much was paid
	PaymentDate    time.Time `json:"payment_date"`    // When the buyer paid
	PaymentMethod  string    `json:"payment_method"`  // e.g., "Cash", "Check", "WireTransfer"
	TransactionRef string    `json:"transaction_ref"` // Bank transaction ID or check number
	CreatedAt      time.Time `json:"created_at"`      // When this record was created
	LastModified   time.Time `json:"last_modified"`   // When last updated
}
