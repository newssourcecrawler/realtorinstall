package models

import "time"

// Property represents a real‐estate listing.
type Property struct {
	ID           int64     `db:"id" json:"id"`
	TenantID     string    `db:"tenant_id" json:"tenantID"` // Which tenant this property belongs to
	Address      string    `db:"address" json:"address"`
	City         string    `db:"city" json:"city"`
	ZIP          string    `db:"zip" json:"zip"`
	ListingDate  time.Time `db:"listing_date" json:"listing_date"`
	CreatedBy    string    `db:"created_by" json:"created_by"` // Username or userID who created
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	ModifiedBy   string    `db:"modified_by" json:"modified_by"` // Username or userID who last modified
	LastModified time.Time `db:"last_modified" json:"last_modified"`
	Deleted      bool      `db:"deleted" json:"deleted"`
}
