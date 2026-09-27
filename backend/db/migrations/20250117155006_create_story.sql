-- Frozen baseline, see 20250117155001_create_source.sql.

-- +goose Up
CREATE TABLE stories (
    id          BIGSERIAL   NOT NULL,
    headline_id BIGINT      NOT NULL,
    summary     TEXT,
    content     TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    CONSTRAINT stories_pkey PRIMARY KEY (id),
    CONSTRAINT stories_headline_id_key UNIQUE (headline_id)
);

-- +goose Down
DROP TABLE IF EXISTS stories CASCADE;
