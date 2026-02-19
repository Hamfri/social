CREATE TABLE IF NOT EXISTS roles(
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    level int NOT NULL DEFAULT 0,
    description TEXT
);

INSERT INTO roles (name, level, description)
VALUES
    ('user', 1, 'can create posts and comments'),
    ('moderator', 2, 'can update other users posts'),
    ('admin', 3, 'can update and delete other users posts');
