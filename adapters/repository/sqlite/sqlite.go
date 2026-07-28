package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

type TaskRepository interface {
	SaveTask(ctx context.Context, name, description string, deadline time.Time) (int64, error)
}

type Storage struct {
	db *sql.DB
}

func New(db_path string) (*Storage, error) {
	const op = "internal.delivery.sqlite"

	db, err := sql.Open("sqlite3", db_path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := goose.Up(db, "./migrations"); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var storage Storage
	storage.db = db
	return &storage, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) SaveTask(ctx context.Context, name, description string, deadline time.Time) (int64, error) {
	const op = "internal.delivery.sqlite.SaveTask"

	stmt, err := s.db.Prepare(`INSERT INTO tasks (name, description, deadline)
	VALUES (?, ?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

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
