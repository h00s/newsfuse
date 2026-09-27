-- +goose Up
-- HCL's scraper is registered under source 11, but the seed migration never created the source.
-- Production databases already have it, added by hand.
INSERT INTO sources (id, name, is_scrapable, topic_id) VALUES (11, 'HCL', true, 4)
ON CONFLICT (id) DO NOTHING;
SELECT setval(pg_get_serial_sequence('sources', 'id'), (SELECT max(id) FROM sources));

-- Summaries are plain text now; older ones were wrapped in a single <p>.
UPDATE stories
SET summary = btrim(substring(summary FROM '^\s*<p>(.*)</p>\s*$'))
WHERE summary ~ '^\s*<p>.*</p>\s*$';

-- +goose Down
-- Nothing to undo: source 11 must stay for its scraper, and plain-text summaries are what the
-- previous code read as well (it only re-wrapped new ones).
SELECT 1;
