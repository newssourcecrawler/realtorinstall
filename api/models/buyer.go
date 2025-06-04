package models

import "time"

// Buyer represents a purchaser or customer.
type Buyer struct {
	ID           int64     `db:"id" json:"id"`
	TenantID     string    `db:"tenant_id" json:"tenantID"`
	FirstName    string    `db:"first_name" json:"first_name"`
	LastName     string    `db:"last_name" json:"last_name"`
	Email        string    `db:"email" json:"email"`
	Phone        string    `db:"phone" json:"phone"`
	CreatedBy    string    `db:"created_by" json:"created_by"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	ModifiedBy   string    `db:"modified_by" json:"modified_by"`
	LastModified time.Time `db:"last_modified" json:"last_modified"`
	Deleted      bool      `db:"deleted" json:"deleted"`
}
