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