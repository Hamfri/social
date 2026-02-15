CREATE TABLE IF NOT EXISTS user_tokens (
    token BYTEA PRIMARY KEY NOT NULL,
    user_id BIGINT NOT NULL,
    expiry TIMESTAMPTZ(0) NOT NULL,
    scope TEXT NOT NULL,

    FOREIGN KEY(user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_tokens_user_id ON user_tokens(user_id);
