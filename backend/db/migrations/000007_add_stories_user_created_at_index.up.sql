-- Add composite index to speed up user-scoped story queries and ordering
CREATE INDEX IF NOT EXISTS idx_stories_user_id_created_at_desc
    ON stories (user_id, created_at DESC);
