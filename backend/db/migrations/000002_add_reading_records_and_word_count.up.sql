-- stories テーブルに reading_records カラムを追加
ALTER TABLE stories ADD COLUMN word_count INTEGER NOT NULL DEFAULT 0 CHECK (word_count >= 0);

-- reading_records テーブルを作成
CREATE TABLE IF NOT EXISTS reading_records (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    story_id INTEGER NOT NULL,
    read_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    word_count INTEGER NOT NULL CHECK (word_count >= 0),

    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_story
        FOREIGN KEY (story_id)
        REFERENCES stories(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_reading_records_user_id ON reading_records(user_id);
CREATE INDEX isx_reading_records_story_id ON reading_records(story_id);