DROP TRIGGER IF EXISTS set_timestamp_stories ON stories;
DROP TRIGGER IF EXISTS set_timestamp_users ON users;
DROP FUNCTION IF EXISTS trigger_set_timestamp();

DROP TABLE IF EXISTS stories;
DROP TABLE IF EXISTS users;