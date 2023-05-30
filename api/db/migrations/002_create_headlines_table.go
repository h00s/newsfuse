package migrations

const CreateHeadlinesTable = `
	BEGIN;

	CREATE TABLE IF NOT EXISTS headlines (
    id BIGSERIAL PRIMARY KEY,
    source_id INT NOT NULL REFERENCES sources(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    url TEXT UNIQUE NOT NULL,
    inserted_at TIMESTAMP WITH TIME ZONE NOT NULL
	);

	UPDATE schema_migrations SET version = 2 WHERE version = 0;

	COMMIT;
`
