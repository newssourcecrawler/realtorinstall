-- migrations/buyers/0001_create_buyers_table.sql

CREATE TABLE IF NOT EXISTS buyers (
  id            SERIAL PRIMARY KEY,
  tenant_id     VARCHAR   NOT NULL,
  name          VARCHAR   NOT NULL,
  contact       VARCHAR,
  created_by    VARCHAR   NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  modified_by   VARCHAR   NOT NULL,
  last_modified TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted       BOOLEAN   NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_buyers_tenant_deleted ON buyers(tenant_id, deleted);
