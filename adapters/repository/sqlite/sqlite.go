package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

type TaskRepository interface {
	SaveTask(ctx context.Context, name, description string, deadline time.Time) (int64, error)
	GetAllTasks(ctx context.Context) (*sql.Rows, error)
	DeleteTaskByID(ctx context.Context, id int) (int, error)
	UpdateTask(ctx context.Context, id int, fileds map[string]interface{}) (int, error)
}

type Storage struct {
	db *sql.DB
}

func New(db_path string) (*Storage, error) {
	const op = "adapter.repo.sqlite"

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
	const op = "adapter.repo.sqlite.SaveTask"

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

func (s *Storage) GetAllTasks(ctx context.Context) (*sql.Rows, error) {
	const op = "adapter.repo.sqlite.GetAllTasks"

	res, err := s.db.QueryContext(ctx, `SELECT id, name, description, deadline FROM tasks`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}

func (s *Storage) DeleteTaskByID(ctx context.Context, id int) (int, error) {
	const op = "adapter.repo.sqlite.DeleteTaskByID"

	_, err := s.db.ExecContext(ctx, `DELETE FROM tasks WHERE tasks.id = ?`, id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil

}

func (s *Storage) UpdateTask(ctx context.Context, id int, fileds map[string]interface{}) (int, error) {
	const op = "adapter.repo.sqlite.UpdateTask"
	if len(fileds) == 0 {
		return 0, fmt.Errorf("%s: no fields to update", op)
	}

	allowedFileds := map[string]bool{
		"name":        true,
		"description": true,
		"deadline":    true,
	}

	var args []string
	arg_id := 1
	var argsValue []interface{}

	for field, value := range fileds {
		if !allowedFileds[field] {
			continue
		}
		args = append(args, fmt.Sprintf("%s = ?", field))
		argsValue = append(argsValue, value)
		arg_id++
	}

	req := fmt.Sprintf("UPDATE tasks SET %s WHERE tasks.id = %d", strings.Join(args, ","), id)
	_, err := s.db.ExecContext(ctx, req, argsValue...)
	if err != nil {
		return 0, fmt.Errorf("%s: %w : %s", op, err, req)
	}
	return id, nil
}
