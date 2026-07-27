package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(db_path string) (*Storage, error) {
	const op = "internal.delivery.sqlite"

	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS tasks (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT,
	deadline DATE );
	CREATE INDEX IF NOT EXISTS id_name ON tasks(name) `)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var storage Storage
	storage.db = db
	return &storage, nil
}

func (s *Storage) SaveTask(ctx context.Context, name, description string, deadline time.Time) (int64, error) {
	const op = "internal.delivery.sqlite.SaveTask"

	stmt, err := s.db.Prepare(`INSERT INTO tasks (name, description, deadline)
	VALUES (?, ?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	res, err := stmt.ExecContext(ctx, name, description, deadline)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
