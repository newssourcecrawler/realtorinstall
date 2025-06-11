-- migrations/sales/0001_create_sales_table.sql

CREATE TABLE IF NOT EXISTS sales (
  id           SERIAL PRIMARY KEY,
  tenant_id    VARCHAR   NOT NULL,
  property_id  INTEGER   NOT NULL REFERENCES properties(id),
  buyer_id     INTEGER   NOT NULL REFERENCES buyers(id),
  sale_date    DATE      NOT NULL,
  sale_price   DOUBLE PRECISION NOT NULL,
  sale_type    VARCHAR   NOT NULL,
  created_by   VARCHAR   NOT NULL,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  modified_by  VARCHAR   NOT NULL,
  last_modified TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted      BOOLEAN   NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_sales_tenant_deleted ON sales(tenant_id, deleted);
