// internal/models/property.go
package models

import "time"

// Property represents a real‐estate listing
type Property struct {
	ID           int64     `json:"id"`
	Address      string    `json:"address"`
	City         string    `json:"city"`
	ZIP          string    `json:"zip"`
	LocationCode string    `json:"location_code"`
	SizeSqFt     float64   `json:"size_sqft"`
	BasePriceUSD float64   `json:"base_price_usd"`
	ListingDate  time.Time `json:"listing_date"`
	LastModified time.Time `json:"last_modified"`
}
