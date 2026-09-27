-- +goose Up
INSERT INTO sources (name, topic_id, is_scrapable) VALUES
    ('klikni.hr', 1, true),
    ('MojPortal.hr', 1, true),
    ('Radio Daruvar', 1, true),
    ('Index.hr', 2, true),
    ('N1Info.hr', 2, true),
    ('Index.hr', 3, true),
    ('N1Info.hr', 3, true),
    ('Hacker News', 4, true),
    ('Bug', 4, true),
    ('Telegram', 2, true);

-- +goose Down
DELETE FROM sources
WHERE name IN ('klikni.hr', 'MojPortal.hr', 'Radio Daruvar', 'Index.hr', 'N1Info.hr', 'Hacker News', 'Bug', 'Telegram');
