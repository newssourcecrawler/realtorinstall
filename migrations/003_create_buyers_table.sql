CREATE TABLE IF NOT EXISTS buyers (
	  id INTEGER PRIMARY KEY AUTOINCREMENT,
	  tenant_id TEXT NOT NULL,
	  first_name TEXT NOT NULL,
	  last_name TEXT NOT NULL,
	  email TEXT NOT NULL,
	  phone TEXT,
	  created_by TEXT NOT NULL,
	  created_at DATETIME NOT NULL,
	  modified_by TEXT NOT NULL,
	  last_modified DATETIME NOT NULL,
	  deleted INTEGER NOT NULL DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_buyers_tenant ON buyers(tenant_id);