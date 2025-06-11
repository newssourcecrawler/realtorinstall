-- 0003_create_roles_permissions.sql

-- Roles master table
CREATE TABLE IF NOT EXISTS roles (
  id          SERIAL PRIMARY KEY,
  name        TEXT    NOT NULL UNIQUE,
  description TEXT
);

-- Permissions master table
CREATE TABLE IF NOT EXISTS permissions (
  id          SERIAL PRIMARY KEY,
  name        TEXT    NOT NULL UNIQUE,
  description TEXT
);

-- Junction: which permissions each role has
CREATE TABLE IF NOT EXISTS role_permissions (
  role_id       INTEGER NOT NULL REFERENCES roles(id),
  permission_id INTEGER NOT NULL REFERENCES permissions(id),
  PRIMARY KEY(role_id, permission_id)
);

-- Junction: which roles each user has
CREATE TABLE IF NOT EXISTS user_roles (
  user_id INTEGER NOT NULL REFERENCES users(id),
  role_id INTEGER NOT NULL REFERENCES roles(id),
  PRIMARY KEY(user_id, role_id)
);

-- Indexes for faster lookups
CREATE INDEX IF NOT EXISTS idx_role_permissions_role ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_user       ON user_roles(user_id);
