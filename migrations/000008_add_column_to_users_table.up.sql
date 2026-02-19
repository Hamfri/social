-- poor database design 
-- use a pivot table instead of adding new columns
ALTER TABLE IF EXISTS users
    ADD COLUMN role_id BIGINT REFERENCES roles(id) DEFAULT 1;

CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);

UPDATE users SET role_id = (SELECT id FROM roles WHERE name = 'user');

ALTER TABLE users ALTER COLUMN role_id DROP DEFAULT;
ALTER TABLE users ALTER COLUMN role_id SET NOT NULL;
