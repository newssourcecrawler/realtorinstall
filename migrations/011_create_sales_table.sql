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