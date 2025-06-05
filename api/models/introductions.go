package models

import "time"

type Introductions struct {
	ID              int64     `db:"id" json:"id"`
	TenantID        string    `db:"tenant_id" json:"tenantID"`
	IntroducerID    int64     `db:"introducer_id" json:"introducerID"`
	IntroducedParty string    `db:"introduced_party" json:"introducedparty"`
	PropertyID      int64     `db:"property_id" json:"propertyID"`
	TransactionID   int64     `db:"transaction_id" json:"transactionID"`
	TransactionType string    `db:"transaction_type" json:"transactiontype"`
	IntroDate       time.Time `db:"intro_date" json:"introdate"`
	AgreedFee       float64   `db:"agreed_fee" json:"agreedfee"`
	FeeType         string    `db:"fee_type" json:"feetype"`
	CreatedBy       string    `db:"created_by" json:"created_by"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	ModifiedBy      string    `db:"modified_by" json:"modified_by"`
	LastModified    time.Time `db:"last_modified" json:"last_modified"`
	Deleted         bool      `db:"deleted" json:"deleted"`
}
