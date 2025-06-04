package models

import "time"

// LocationPricing stores price guidance per square foot by ZIP code.
type LocationPricing struct {
	ID            int64     `db:"id" json:"id"`
	TenantID      string    `db:"tenant_id" json:"tenantID"`
	ZipCode       string    `db:"zip_code" json:"zip_code"` // Postal code, e.g. "78701"
	City          string    `db:"city" json:"city"`         // (Optional) city name
	PricePerSqFt  float64   `db:"price_per_sqft" json:"price_per_sqft"`
	EffectiveDate time.Time `db:"effective_date" json:"effective_date"`
	CreatedBy     string    `db:"created_by" json:"created_by"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	ModifiedBy    string    `db:"modified_by" json:"modified_by"`
	LastModified  time.Time `db:"last_modified" json:"last_modified"`
	Deleted       bool      `db:"deleted" json:"deleted"`
}
