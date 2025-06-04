package models

// Buyer represents a purchaser or customer.
type Buyer struct {
	ID           int64  `db:"id" json:"id"`
	Name         string `db:"name" json:"name"`
	Email        string `db:"email" json:"email"`
	CreatedAt    string `db:"created_at" json:"created_at"`
	LastModified string `db:"last_modified" json:"last_modified"`
	Deleted      bool   `db:"deleted" json:"deleted"`
}
