-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = current_timestamp;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

ALTER TABLE topics ALTER COLUMN name SET NOT NULL;

UPDATE sources SET is_scrapable = false WHERE is_scrapable IS NULL;
ALTER TABLE sources
    ALTER COLUMN name SET NOT NULL,
    ALTER COLUMN is_scrapable SET DEFAULT false,
    ALTER COLUMN is_scrapable SET NOT NULL,
    ADD CONSTRAINT sources_topic_id_fkey FOREIGN KEY (topic_id) REFERENCES topics (id) ON DELETE RESTRICT;
CREATE INDEX sources_topic_id_idx ON sources (topic_id);

-- A row without a title, url or source can't be shown or linked; there are none today.
DELETE FROM headlines WHERE title IS NULL OR url IS NULL OR source_id IS NULL;
ALTER TABLE headlines
    DROP COLUMN IF EXISTS story,
    ALTER COLUMN title SET NOT NULL,
    ALTER COLUMN url SET NOT NULL,
    ALTER COLUMN source_id SET NOT NULL;
CREATE TRIGGER headlines_set_updated_at
    BEFORE UPDATE ON headlines
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Stories whose headline is gone can never be reached, and they would block the foreign key.
DELETE FROM stories WHERE NOT EXISTS (SELECT 1 FROM headlines WHERE headlines.id = stories.headline_id);
UPDATE stories SET summary = '' WHERE summary IS NULL;
UPDATE stories SET content = '' WHERE content IS NULL;
ALTER TABLE stories
    ALTER COLUMN summary SET DEFAULT '',
    ALTER COLUMN summary SET NOT NULL,
    ALTER COLUMN content SET DEFAULT '',
    ALTER COLUMN content SET NOT NULL,
    ADD CONSTRAINT stories_headline_id_fkey FOREIGN KEY (headline_id) REFERENCES headlines (id) ON DELETE CASCADE;
CREATE TRIGGER stories_set_updated_at
    BEFORE UPDATE ON stories
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- +goose Down
DROP TRIGGER IF EXISTS stories_set_updated_at ON stories;
ALTER TABLE stories
    DROP CONSTRAINT IF EXISTS stories_headline_id_fkey,
    ALTER COLUMN content DROP NOT NULL,
    ALTER COLUMN content DROP DEFAULT,
    ALTER COLUMN summary DROP NOT NULL,
    ALTER COLUMN summary DROP DEFAULT;

DROP TRIGGER IF EXISTS headlines_set_updated_at ON headlines;
ALTER TABLE headlines
    ALTER COLUMN source_id DROP NOT NULL,
    ALTER COLUMN url DROP NOT NULL,
    ALTER COLUMN title DROP NOT NULL,
    ADD COLUMN IF NOT EXISTS story JSONB;

DROP INDEX IF EXISTS sources_topic_id_idx;
ALTER TABLE sources
    DROP CONSTRAINT IF EXISTS sources_topic_id_fkey,
    ALTER COLUMN is_scrapable DROP NOT NULL,
    ALTER COLUMN is_scrapable DROP DEFAULT,
    ALTER COLUMN name DROP NOT NULL;

ALTER TABLE topics ALTER COLUMN name DROP NOT NULL;

DROP FUNCTION IF EXISTS set_updated_at();
