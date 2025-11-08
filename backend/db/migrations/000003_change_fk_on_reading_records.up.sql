ALTER TABLE reading_records DROP CONSTRAINT fk_story;

ALTER TABLE reading_records ALTER COLUMN story_id DROP NOT NULL;

ALTER TABLE reading_records ADD CONSTRAINT fk_story
    FOREIGN KEY (story_id)
    REFERENCES stories(id)
    ON DELETE SET NULL;