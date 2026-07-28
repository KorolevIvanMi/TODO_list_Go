-- +goose Up
CREATE TABLE IF NOT EXISTS tasks (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT,
	deadline DATE );
	CREATE INDEX IF NOT EXISTS id_name ON tasks(name);

-- +goose Down
DROP TABLE tasks;