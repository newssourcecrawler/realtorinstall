-- migrations/commissions/0001_create_commissions_table.sql

CREATE TABLE IF NOT EXISTS commissions (
  id                SERIAL PRIMARY KEY,
  tenant_id         VARCHAR   NOT NULL,
  transaction_type  VARCHAR   NOT NULL,  -- \"sale\", \"letting\", \"introduction\"
  transaction_id    INTEGER   NOT NULL,
  beneficiary_id    INTEGER   NOT NULL REFERENCES buyers(id),
  commission_type   VARCHAR   NOT NULL,  -- \"percentage\" or \"fixed\"
  rate_or_amount    DOUBLE PRECISION NOT NULL,
  calculated_amount DOUBLE PRECISION NOT NULL,
  memo              TEXT,
  created_by        VARCHAR   NOT NULL,
  created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  modified_by       VARCHAR   NOT NULL,
  last_modified     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted           BOOLEAN   NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_commissions_tenant_deleted ON commissions(tenant_id, deleted);
CREATE INDEX idx_commissions_txn      ON commissions(transaction_type, transaction_id);
CREATE INDEX idx_commissions_benef    ON commissions(beneficiary_id);
