package models

import "time"

type Sales struct {
	ID           int64     `db:"id" json:"id"`
	TenantID     string    `db:"tenant_id" json:"tenantID"`
	PropertyID   int64     `db:"property_id" json:"propertyID"`
	BuyerID      int64     `db:"buyer_id" json:"buyerID"`
	SalePrice    float64   `db:"sale_price" json:"saleprice"`
	SaleDate     time.Time `db:"sale_date" json:"saledate"`
	SaleType     string    `db:"sale_type" json:"saletype"`
	CreatedBy    string    `db:"created_by" json:"created_by"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	ModifiedBy   string    `db:"modified_by" json:"modified_by"`
	LastModified time.Time `db:"last_modified" json:"last_modified"`
	Deleted      bool      `db:"deleted" json:"deleted"`
}

type MonthlySummary struct {
	SalesMonth time.Month `json:"sales_month"`
	SalesYear  int        `json:"sales_year"`
	TotalSales float64    `json:"total_sales"`
}
