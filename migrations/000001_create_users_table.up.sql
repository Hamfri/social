CREATE EXTENSION IF NOT EXISTS citext; -- dataType to support case insensitive columns like emails

CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    email CITEXT UNIQUE NOT NULL,
    password BYTEA NOT NULL,
    created_at TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ(0) NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
