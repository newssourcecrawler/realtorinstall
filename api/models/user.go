package models

import "time"

// User represents a system account (salesperson, admin, etc.)
type User struct {
	ID           int64     `json:"id"`            // Primary key
	UserName     string    `json:"username"`      // unique login
	PasswordHash string    `json:"-"`             // bcrypt or argon2 hash
	FirstName    string    `json:"first_name"`    //  first name
	LastName     string    `json:"last_name"`     //  last name
	Role         string    `json:"role"`          // e.g. "admin", "sales"
	Email        string    `json:"email"`         // Contact email
	Phone        string    `json:"phone"`         // Contact phone number
	CreatedAt    time.Time `json:"created_at"`    // When this record was created
	LastModified time.Time `json:"last_modified"` // When last updated
	CreatedBy    string    `json:"created_by"`
	ModifiedBy   string    `json:"modified_by"`
	Deleted      bool      `json:"deleted"` // Soft‐delete flag
}
