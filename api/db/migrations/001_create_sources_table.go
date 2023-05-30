package migrations

const CreateSourcesTable = `
	BEGIN;

	CREATE TABLE IF NOT EXISTS sources (
		id SERIAL PRIMARY KEY,
		name VARCHAR(80) NOT NULL
	);

	UPDATE schema_migrations SET version = 1 WHERE version = 0;

	COMMIT;
`
