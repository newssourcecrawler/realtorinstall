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