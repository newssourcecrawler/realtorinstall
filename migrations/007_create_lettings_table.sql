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