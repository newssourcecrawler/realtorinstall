CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  tenant_id VARCHAR NOT NULL,
  username VARCHAR NOT NULL,
  password_hash VARCHAR NOT NULL,
  first_name VARCHAR NOT NULL,
  last_name VARCHAR NOT NULL,
  role VARCHAR NOT NULL,
  email VARCHAR NOT NULL,
  phone VARCHAR,
  created_by VARCHAR NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  modified_by VARCHAR NOT NULL,
  last_modified TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX idx_users_tenant_deleted ON users(tenant_id, deleted);
