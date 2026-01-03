-- Add composite indexes to speed up reading_records queries
CREATE INDEX IF NOT EXISTS idx_reading_records_user_id_read_at
    ON reading_records (user_id, read_at);

CREATE INDEX IF NOT EXISTS idx_reading_records_user_story_read_at_desc
    ON reading_records (user_id, story_id, read_at DESC);
