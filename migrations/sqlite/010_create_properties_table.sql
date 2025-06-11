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