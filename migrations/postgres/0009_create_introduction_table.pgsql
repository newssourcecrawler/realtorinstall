-- migrations/introductions/0001_create_introductions_table.sql

CREATE TABLE IF NOT EXISTS introductions (
  id                SERIAL PRIMARY KEY,
  tenant_id         VARCHAR   NOT NULL,
  introducer_id     INTEGER   NOT NULL REFERENCES buyers(id),
  introduced_party  VARCHAR   NOT NULL,
  property_id       INTEGER   NOT NULL REFERENCES properties(id),
  transaction_type  VARCHAR,  -- optionally \"sale\" or \"letting\"
  transaction_id    INTEGER,
  agreed_fee        DOUBLE PRECISION,
  created_by        VARCHAR   NOT NULL,
  created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  modified_by       VARCHAR   NOT NULL,
  last_modified     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted           BOOLEAN   NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_intros_tenant_deleted ON introductions(tenant_id, deleted);
