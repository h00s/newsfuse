-- Frozen baseline, see 20250117155001_create_source.sql. The story column is what Bun made of an
-- untagged relation field; 20260927120000_harden_schema.sql drops it.

-- +goose Up
CREATE TABLE headlines (
    id           BIGSERIAL   NOT NULL,
    title        VARCHAR,
    url          VARCHAR,
    source_id    BIGINT,
    story        JSONB,
    published_at TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    CONSTRAINT headlines_pkey PRIMARY KEY (id),
    CONSTRAINT headlines_url_key UNIQUE (url),
    CONSTRAINT headlines_source_id_fkey FOREIGN KEY (source_id)
        REFERENCES sources (id) ON DELETE CASCADE
);

CREATE INDEX headlines_source_id_idx ON headlines (source_id);
CREATE INDEX headlines_published_at_idx ON headlines (published_at);

-- +goose Down
DROP TABLE IF EXISTS headlines CASCADE;
