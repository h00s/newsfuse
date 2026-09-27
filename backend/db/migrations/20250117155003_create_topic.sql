-- Frozen baseline, see 20250117155001_create_source.sql.

-- +goose Up
CREATE TABLE topics (
    id   BIGSERIAL NOT NULL,
    name VARCHAR,
    CONSTRAINT topics_pkey PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE IF EXISTS topics CASCADE;
