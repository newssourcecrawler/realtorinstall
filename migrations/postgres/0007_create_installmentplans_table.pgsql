-- migrations/plan/0001_create_installment_plans_table.sql

CREATE TABLE IF NOT EXISTS installment_plans (
  id            SERIAL PRIMARY KEY,
  tenant_id     VARCHAR   NOT NULL,
  property_id   INTEGER   NOT NULL REFERENCES properties(id),
  buyer_id      INTEGER   NOT NULL REFERENCES buyers(id),
  total_amount  DOUBLE PRECISION NOT NULL,
  created_by    VARCHAR   NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  modified_by   VARCHAR   NOT NULL,
  last_modified TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted       BOOLEAN   NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_plans_tenant_deleted ON installment_plans(tenant_id, deleted);
