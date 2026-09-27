-- Frozen baseline: reproduces the table Bun created from the original model, so databases that
-- already ran the Go version of this migration and fresh ones end up with the same schema.

-- +goose Up
CREATE TABLE sources (
    id           BIGSERIAL NOT NULL,
    name         VARCHAR,
    is_scrapable BOOLEAN,
    topic_id     BIGINT    NOT NULL,
    CONSTRAINT sources_pkey PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE IF EXISTS sources CASCADE;
