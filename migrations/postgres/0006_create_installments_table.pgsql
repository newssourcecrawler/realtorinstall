-- migrations/installments/0001_create_installments_table.sql

CREATE TABLE IF NOT EXISTS installments (
  id             SERIAL PRIMARY KEY,
  tenant_id      VARCHAR   NOT NULL,
  plan_id        INTEGER   NOT NULL REFERENCES installment_plans(id),
  due_date       DATE      NOT NULL,
  amount_due     DOUBLE PRECISION NOT NULL,
  amount_paid    DOUBLE PRECISION NOT NULL DEFAULT 0,
  created_by     VARCHAR   NOT NULL,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  modified_by    VARCHAR   NOT NULL,
  last_modified  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted        BOOLEAN   NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_installments_tenant_deleted ON installments(tenant_id, deleted);
