package models

import "time"

type Lettings struct {
	ID           int64     `db:"id" json:"id"`
	TenantID     string    `db:"tenant_id" json:"tenantID"`
	PropertyID   int64     `db:"property_id" json:"propertyID"`
	CycleID      int64     `db:"cycle_id" json:"cycleID"`
	TenantUserID int64     `db:"tenant_user_id" json:"tenant_userID"`
	RentAmount   float64   `db:"rent_amount" json:"rentamount"`
	RentTerm     int64     `db:"rent_term" json:"rentterm"`
	RentCycle    string    `db:"rent_cycle" json:"rentcycle"`
	Memo         string    `db:"memo" json:"memo"`
	StartDate    time.Time `db:"start_date" json:"startdate"`
	EndDate      time.Time `db:"end_date" json:"enddate"`
	CreatedBy    string    `db:"created_by" json:"created_by"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	ModifiedBy   string    `db:"modified_by" json:"modified_by"`
	LastModified time.Time `db:"last_modified" json:"last_modified"`
	Deleted      bool      `db:"deleted" json:"deleted"`
}

type RentRoll struct {
	PropertyID int64   `json:"property_id"`
	TotalRent  float64 `json:"total_rent"`
}
