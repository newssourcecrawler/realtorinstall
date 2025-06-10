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