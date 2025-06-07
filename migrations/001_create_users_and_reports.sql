-- 001_create_users_and_reports.sql

-- 1. Users table
CREATE TABLE IF NOT EXISTS users (
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  tenant_id     TEXT    NOT NULL,
  username      TEXT    NOT NULL UNIQUE,
  password_hash TEXT    NOT NULL,
  first_name    TEXT    NOT NULL,
  last_name     TEXT    NOT NULL,
  role          TEXT    NOT NULL,
  email         TEXT    NOT NULL,
  phone         TEXT,
  created_by    TEXT    NOT NULL,
  created_at    DATETIME NOT NULL,
  modified_by   TEXT    NOT NULL,
  last_modified DATETIME NOT NULL,
  deleted       INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_users_tenant ON users(tenant_id);

-- 2. Commissions table
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
CREATE INDEX IF NOT EXISTS idx_commissions_tenant ON commissions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_commissions_txn     ON commissions(transaction_type, transaction_id);
CREATE INDEX IF NOT EXISTS idx_commissions_benef   ON commissions(beneficiary_id);

-- 3. Commission‐by‐beneficiary view
CREATE VIEW IF NOT EXISTS view_commission_by_beneficiary AS
SELECT 
  beneficiary_id,
  SUM(calculated_amount) AS total_commission
FROM commissions
WHERE deleted = 0
GROUP BY beneficiary_id;
