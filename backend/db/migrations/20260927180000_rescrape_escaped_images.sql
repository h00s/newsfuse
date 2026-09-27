-- +goose Up
-- For a few hours text-mode scraping escaped lazy-loaded images, so these stories show a literal
-- "<img …>" instead of the picture. They are scraped again on their next read; a summary, if any,
-- is generated again on request.
DELETE FROM stories WHERE content LIKE '%&lt;img%';

-- +goose Down
-- Nothing to restore: the deleted rows were a cache of the news site.
SELECT 1;
