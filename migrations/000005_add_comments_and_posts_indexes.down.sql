DROP INDEX IF EXISTS idx_posts_title;
DROP INDEX IF EXISTS idx_posts_tags;
DROP INDEX IF EXISTS idx_comments_comment;

DROP EXTENSION IF EXISTS pg_trgm;
