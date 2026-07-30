package updatetask

import (
	"context"
	"fmt"
	"time"

	"github.com/KorolevIvanMi/TODO_list_Go/adapters/repository/sqlite"
	"github.com/KorolevIvanMi/TODO_list_Go/internal/domain"
)

type UseCase struct {
	repo sqlite.TaskRepository
}

func New(repo sqlite.TaskRepository) *UseCase {
	uc := UseCase{repo: repo}
	return &uc
}

func (uc *UseCase) UpdateTask(ctx context.Context, id int, name, description *string, deadline *time.Time) (int, error) {
	const op = "internal.usecase.task.usecase.UpdateTask"
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}
	args, err := domain.UpdateTask(name, description, deadline)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	resId, err := uc.repo.UpdateTask(ctx, id, args)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return resId, nil
}
