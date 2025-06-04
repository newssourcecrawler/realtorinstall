// internal/models/property.go
package models

// Property represents a real‐estate listing
type Property struct {
	ID           int64  `db:"id" json:"id"`
	Address      string `db:"address" json:"address"`
	City         string `db:"city" json:"city"`
	ZIP          string `db:"zip" json:"zip"`
	ListingDate  string `db:"listing_date" json:"listing_date"`
	CreatedAt    string `db:"created_at" json:"created_at"`
	LastModified string `db:"last_modified" json:"last_modified"`
	Deleted      bool   `db:"deleted" json:"deleted"`
}
