ALTER TABLE reading_records DROP CONSTRAINT fk_story;

DELETE FROM reading_records WHERE story_id IS NULL;

ALTER TABLE reading_records ALTER COLUMN story_id SET NOT NULL;

ALTER TABLE reading_records ADD CONSTRAINT fk_story
    FOREIGN KEY (story_id)
    REFERENCES stories(id)
    ON DELETE CASCADE;