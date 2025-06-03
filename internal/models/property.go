// internal/models/property.go
package models

// Property represents a real‐estate listing
type Property struct {
	ID           int64  `json:"id"`
	Address      string `json:"address"`
	City         string `json:"city"`
	ZIP          string `json:"zip"`
	ListingDate  string `json:"listing_date"`  // use RFC3339 string instead of time.Time
	CreatedAt    string `json:"created_at"`    // same here
	LastModified string `json:"last_modified"` // same here
}
