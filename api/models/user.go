package models

import "time"

// User represents a system account (salesperson, admin, etc.)
type User struct {
	ID           int64     `db:"id" json:"id"`                       // Primary key
	TenantID     string    `db:"tenant_id" json:"tenantID"`          // Which tenant this user belongs to
	UserName     string    `db:"username" json:"username"`           // Unique login name
	PasswordHash string    `db:"password_hash" json:"-"`             // bcrypt/argon2 hash (never JSON‐export this)
	FirstName    string    `db:"first_name" json:"first_name"`       // First name
	LastName     string    `db:"last_name" json:"last_name"`         // Last name
	Role         string    `db:"role" json:"role"`                   // e.g. "admin", "sales"
	Email        string    `db:"email" json:"email"`                 // Contact email
	Phone        string    `db:"phone" json:"phone"`                 // Contact phone number
	CreatedBy    string    `db:"created_by" json:"created_by"`       // Who created this user
	CreatedAt    time.Time `db:"created_at" json:"created_at"`       // When this record was created
	ModifiedBy   string    `db:"modified_by" json:"modified_by"`     // Who last modified
	LastModified time.Time `db:"last_modified" json:"last_modified"` // When last updated
	Deleted      bool      `db:"deleted" json:"deleted"`             // Soft‐delete flag
}
