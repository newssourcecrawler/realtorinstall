-- migrations/property/0001_create_properties_table.sql

CREATE TABLE IF NOT EXISTS properties (
  id             SERIAL PRIMARY KEY,
  tenant_id      VARCHAR   NOT NULL,
  location_code  VARCHAR   NOT NULL,
  size_sq_ft     DOUBLE PRECISION NOT NULL,
  base_price_usd DOUBLE PRECISION NOT NULL,
  listing_date   TIMESTAMPTZ NOT NULL,
  created_by     VARCHAR   NOT NULL,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  modified_by    VARCHAR   NOT NULL,
  last_modified  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted        BOOLEAN   NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_properties_tenant_deleted ON properties(tenant_id, deleted);
