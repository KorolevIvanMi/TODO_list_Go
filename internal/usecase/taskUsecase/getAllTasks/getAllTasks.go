package getalltasks

import (
	"context"
	"fmt"
	"time"

	"github.com/KorolevIvanMi/TODO_list_Go/adapters/repository/sqlite"
)

type UseCase struct {
	repo sqlite.TaskRepository
}

func New(repo sqlite.TaskRepository) *UseCase {
	uc := UseCase{repo: repo}
	return &uc
}

type GetedTask struct {
	Id          int       `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Deadline    time.Time `db:"deadline"`
}

func (uc *UseCase) GetAllTasks(ctx context.Context) (*[]GetedTask, error) {
	const op = "internal.usecase.task.usecase.CreateTask"

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	tasks, err := uc.repo.GetAllTasks(ctx)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer tasks.Close()
	var result []GetedTask
	for tasks.Next() {
		var task GetedTask
		err = tasks.Scan(&task.Id, &task.Name, &task.Description, &task.Deadline)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		result = append(result, task)
	}
	if err = tasks.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &result, nil
}
