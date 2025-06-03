package models

import "time"

// LocationPricing stores price guidance per square foot by ZIP code.
type LocationPricing struct {
	ID            int64     `json:"id"`             // Primary key
	ZipCode       string    `json:"zip_code"`       // Postal code, e.g. "78701"
	City          string    `json:"city"`           // Optional city name for reference
	PricePerSqFt  float64   `json:"price_per_sqft"` // USD per square foot
	EffectiveDate time.Time `json:"effective_date"` // When this rate takes effect
	CreatedAt     time.Time `json:"created_at"`     // When this record was created
	LastModified  time.Time `json:"last_modified"`  // When last updated
}
