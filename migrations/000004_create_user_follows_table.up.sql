CREATE TABLE IF NOT EXISTS user_follows (
    followed_id BIGINT NOT NULL, -- main user
    follower_id BIGINT NOT NULL, -- other user accounts following the main user
    created_at TIMESTAMPTZ(0) NOT NULL DEFAULT NOW(),

    CONSTRAINT user_follows_no_self_follow
        CHECK (follower_id <> followed_id),

    FOREIGN KEY (followed_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    
    PRIMARY KEY(followed_id, follower_id)
);

CREATE INDEX IF NOT EXISTS idx_user_follows_followed_id ON user_follows(followed_id);
CREATE INDEX IF NOT EXISTS idx_user_follows_follower_id ON user_follows(follower_id);
