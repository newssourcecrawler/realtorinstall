-- migrations/locationpricing/0001_create_location_pricing_table.sql

CREATE TABLE IF NOT EXISTS location_pricing (
  id            SERIAL PRIMARY KEY,
  tenant_id     VARCHAR   NOT NULL,
  location_code VARCHAR   NOT NULL,
  price_per_sqft DOUBLE PRECISION NOT NULL,
  created_by    VARCHAR   NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  modified_by   VARCHAR   NOT NULL,
  last_modified TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted       BOOLEAN   NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_pricing_tenant_deleted ON location_pricing(tenant_id, deleted);
