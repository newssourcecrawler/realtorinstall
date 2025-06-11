-- migrations/payments/0001_create_payments_table.sql

CREATE TABLE IF NOT EXISTS payments (
  id             SERIAL PRIMARY KEY,
  tenant_id      VARCHAR   NOT NULL,
  installment_id INTEGER   NOT NULL REFERENCES installments(id),
  payment_date   DATE      NOT NULL,
  amount         DOUBLE PRECISION NOT NULL,
  created_by     VARCHAR   NOT NULL,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  modified_by    VARCHAR   NOT NULL,
  last_modified  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted        BOOLEAN   NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_payments_tenant_deleted ON payments(tenant_id, deleted);
