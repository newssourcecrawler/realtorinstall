-- migrations/lettings/0001_create_lettings_table.sql

CREATE TABLE IF NOT EXISTS lettings (
  id             SERIAL PRIMARY KEY,
  tenant_id      VARCHAR   NOT NULL,
  property_id    INTEGER   NOT NULL REFERENCES properties(id),
  tenant_user_id INTEGER   NOT NULL REFERENCES buyers(id),
  rent_amount    DOUBLE PRECISION NOT NULL,
  start_date     DATE      NOT NULL,
  end_date       DATE,
  created_by     VARCHAR   NOT NULL,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  modified_by    VARCHAR   NOT NULL,
  last_modified  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted        BOOLEAN   NOT NULL DEFAULT FALSE
);

CREATE INDEX idx_lettings_tenant_deleted ON lettings(tenant_id, deleted);
