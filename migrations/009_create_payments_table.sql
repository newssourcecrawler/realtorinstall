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