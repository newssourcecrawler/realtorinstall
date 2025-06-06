package dbmigrations

import (
	"database/sql"
	"fmt"
)

var migrations = []struct {
	name string
	sql  string
}{
	{
		name: "create_users_table",
		sql: `
CREATE TABLE IF NOT EXISTS users (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  tenant_id     TEXT    NOT NULL,
  username      TEXT    NOT NULL UNIQUE,
  password_hash TEXT    NOT NULL,
  first_name    TEXT    NOT NULL,
  last_name     TEXT    NOT NULL,
  email         TEXT    NOT NULL,
  phone         TEXT,
  created_by    TEXT    NOT NULL,
  created_at    DATETIME NOT NULL,
  modified_by   TEXT    NOT NULL,
  last_modified DATETIME NOT NULL,
  deleted       INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_users_tenant ON users(tenant_id);
`,
	},
	{
		name: "create_roles_table",
		sql: `
CREATE TABLE IF NOT EXISTS roles (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        TEXT NOT NULL UNIQUE,
  description TEXT
);
`,
	},
	{
		name: "create_permissions_table",
		sql: `
CREATE TABLE IF NOT EXISTS permissions (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        TEXT NOT NULL UNIQUE,
  description TEXT
);
`,
	},
	{
		name: "create_role_permissions_table",
		sql: `
CREATE TABLE IF NOT EXISTS role_permissions (
  role_id       INTEGER NOT NULL REFERENCES roles(id),
  permission_id INTEGER NOT NULL REFERENCES permissions(id),
  PRIMARY KEY(role_id, permission_id)
);
`,
	},
	{
		name: "create_user_roles_table",
		sql: `
CREATE TABLE IF NOT EXISTS user_roles (
  user_id INTEGER NOT NULL REFERENCES users(id),
  role_id INTEGER NOT NULL REFERENCES roles(id),
  PRIMARY KEY(user_id, role_id)
);
`,
	}, {
		name: "create_sales_table",
		sql: `
	CREATE TABLE IF NOT EXISTS sales (
	  id            INTEGER PRIMARY KEY AUTOINCREMENT,
	  tenant_id     TEXT    NOT NULL,
	  property_id   INTEGER NOT NULL,
	  buyer_id      INTEGER NOT NULL,
	  sale_price    REAL    NOT NULL,
	  sale_date     DATETIME NOT NULL,
	  sale_type     TEXT    NOT NULL,
	  created_by    TEXT    NOT NULL,
	  created_at    DATETIME NOT NULL,
	  modified_by   TEXT    NOT NULL,
	  last_modified DATETIME NOT NULL,
	  deleted       INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_sales_tenant     ON sales(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_sales_property   ON sales(property_id);
	CREATE INDEX IF NOT EXISTS idx_sales_buyer      ON sales(buyer_id);
	`,
	}, {
		name: "create_properties_table",
		sql: `
	CREATE TABLE IF NOT EXISTS properties (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  tenant_id TEXT NOT NULL,
	  address TEXT NOT NULL,
	  city TEXT NOT NULL,
	  zip TEXT NOT NULL,
	  listing_date DATETIME NOT NULL,
	  created_by TEXT NOT NULL,
	  created_at DATETIME NOT NULL,
	  modified_by TEXT NOT NULL,
	  last_modified DATETIME NOT NULL,
	  deleted INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_properties_tenant ON properties(tenant_id);
	`,
	}, {
		name: "create_payments_table",
		sql: `
	CREATE TABLE IF NOT EXISTS payments (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  tenant_id TEXT NOT NULL,
	  installment_id INTEGER NOT NULL,
	  amount_paid REAL NOT NULL,
	  payment_date DATETIME NOT NULL,
	  payment_method TEXT NOT NULL,
	  transaction_ref TEXT,
	  created_by TEXT NOT NULL,
	  created_at DATETIME NOT NULL,
	  modified_by TEXT NOT NULL,
	  last_modified DATETIME NOT NULL,
	  deleted INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_payments_tenant ON payments(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_payments_installment ON payments(installment_id);
	`,
	}, {
		name: "create_pricing_table",
		sql: `
	CREATE TABLE IF NOT EXISTS location_pricing (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  tenant_id TEXT NOT NULL,
	  zip_code TEXT NOT NULL,
	  city TEXT NOT NULL,
	  price_per_sqft REAL NOT NULL,
	  effective_date DATETIME NOT NULL,
	  created_by TEXT NOT NULL,
	  created_at DATETIME NOT NULL,
	  modified_by TEXT NOT NULL,
	  last_modified DATETIME NOT NULL,
	  deleted INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_locationpricing_tenant ON location_pricing(tenant_id);
	`,
	}, {
		name: "create_lettings_table",
		sql: `
	CREATE TABLE IF NOT EXISTS lettings (
	  id             INTEGER PRIMARY KEY AUTOINCREMENT,
	  tenant_id      TEXT    NOT NULL,
	  property_id    INTEGER NOT NULL,
	  tenant_user_id INTEGER NOT NULL,
	  rent_amount    REAL    NOT NULL,
	  rent_term      INTEGER NOT NULL,
	  rent_cycle     TEXT    NOT NULL,
	  memo           TEXT,
	  start_date     DATETIME NOT NULL,
	  end_date       DATETIME,
	  created_by     TEXT    NOT NULL,
	  created_at     DATETIME NOT NULL,
	  modified_by    TEXT    NOT NULL,
	  last_modified  DATETIME NOT NULL,
	  deleted        INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_lettings_tenant    ON lettings(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_lettings_property  ON lettings(property_id);
	CREATE INDEX IF NOT EXISTS idx_lettings_tenantuser ON lettings(tenant_user_id);
	`,
	}, {
		name: "create_introductions_table",
		sql: `
	CREATE TABLE IF NOT EXISTS introductions (
	  id                INTEGER PRIMARY KEY AUTOINCREMENT,
	  tenant_id         TEXT    NOT NULL,
	  introducer_id     INTEGER NOT NULL,
	  introduced_party  TEXT    NOT NULL,
	  property_id       INTEGER NOT NULL,
	  transaction_id    INTEGER,
	  transaction_type	TEXT NOT NULL,
	  intro_date        DATETIME NOT NULL,
	  agreed_fee        REAL    NOT NULL,
	  fee_type          TEXT    NOT NULL,
	  created_by        TEXT    NOT NULL,
	  created_at        DATETIME NOT NULL,
	  modified_by       TEXT    NOT NULL,
	  last_modified     DATETIME NOT NULL,
	  deleted           INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_introductions_tenant ON introductions(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_introductions_introducer ON introductions(introducer_id);
	CREATE INDEX IF NOT EXISTS idx_introductions_property ON introductions(property_id);
	`,
	}, {
		name: "create_installmentplans_table",
		sql: `
	CREATE TABLE IF NOT EXISTS installment_plans (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  tenant_id TEXT NOT NULL,
	  property_id INTEGER NOT NULL,
	  buyer_id INTEGER NOT NULL,
	  total_price REAL NOT NULL,
	  down_payment REAL NOT NULL,
	  num_installments INTEGER NOT NULL,
	  frequency TEXT NOT NULL,
	  first_installment DATETIME NOT NULL,
	  interest_rate REAL NOT NULL,
	  created_by TEXT NOT NULL,
	  created_at DATETIME NOT NULL,
	  modified_by TEXT NOT NULL,
	  last_modified DATETIME NOT NULL,
	  deleted INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_installmentplans_tenant ON installment_plans(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_installmentplans_property ON installment_plans(property_id);
	`,
	},
	{
		name: "create_installments_table",
		sql: `
	CREATE TABLE IF NOT EXISTS installments (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  tenant_id TEXT NOT NULL,
	  plan_id INTEGER NOT NULL,
	  sequence_number INTEGER NOT NULL,
	  due_date DATETIME NOT NULL,
	  amount_due REAL NOT NULL,
	  amount_paid REAL NOT NULL,
	  status TEXT NOT NULL,
	  late_fee REAL NOT NULL,
	  paid_date DATETIME,
	  created_by TEXT NOT NULL,
	  created_at DATETIME NOT NULL,
	  modified_by TEXT NOT NULL,
	  last_modified DATETIME NOT NULL,
	  deleted INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_installments_tenant ON installments(tenant_id);
	CREATE INDEX IF NOT EXISTS idx_installments_plan ON installments(plan_id);
	`,
	},
	{
		name: "create_commissions_table",
		sql: `
CREATE TABLE IF NOT EXISTS commissions (
  id                INTEGER PRIMARY KEY AUTOINCREMENT,
  tenant_id         TEXT    NOT NULL,
  transaction_type  TEXT    NOT NULL,
  transaction_id    INTEGER NOT NULL,
  beneficiary_id    INTEGER NOT NULL,
  commission_type   TEXT    NOT NULL,
  rate_or_amount    REAL    NOT NULL,
  calculated_amount REAL    NOT NULL,
  memo              TEXT,
  created_by        TEXT    NOT NULL,
  created_at        DATETIME NOT NULL,
  modified_by       TEXT    NOT NULL,
  last_modified     DATETIME NOT NULL,
  deleted           INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_comm_tenant ON commissions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_comm_txn    ON commissions(transaction_type, transaction_id);
CREATE INDEX IF NOT EXISTS idx_comm_benef  ON commissions(beneficiary_id);
`,
	},
}

func ApplyMigrations(db *sql.DB) error {
	for _, m := range migrations {
		if _, err := db.Exec(m.sql); err != nil {
			return fmt.Errorf("migration %q failed: %w", m.name, err)
		}
	}
	return nil
}
