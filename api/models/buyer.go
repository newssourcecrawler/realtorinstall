package models

import "time"

// Buyer represents a purchaser or customer.
type Buyer struct {
	ID           int64     `json:"id"`            // Primary key
	FirstName    string    `json:"first_name"`    // Buyer’s first name
	LastName     string    `json:"last_name"`     // Buyer’s last name
	Email        string    `json:"email"`         // Contact email
	Phone        string    `json:"phone"`         // Contact phone number
	CreatedAt    time.Time `json:"created_at"`    // When this record was created
	LastModified time.Time `json:"last_modified"` // When last updated
	Deleted      bool      `json:"deleted"`       // Soft‐delete flag
}
