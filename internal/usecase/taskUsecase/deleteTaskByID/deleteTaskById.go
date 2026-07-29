package deletetaskbyid

import (
	"context"
	"fmt"

	"github.com/KorolevIvanMi/TODO_list_Go/adapters/repository/sqlite"
)

type UseCase struct {
	repo sqlite.TaskRepository
}

func New(repo sqlite.TaskRepository) *UseCase {
	uc := UseCase{repo: repo}
	return &uc
}

func (uc *UseCase) DeleteTaskByID(ctx context.Context, id int) (int, error) {
	const op = "internal.usecase.task.usecase.DeleteTaskById"
	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
	}

	idx, err := uc.repo.DeleteTaskByID(ctx, id)
	if err != nil {
		return -0, fmt.Errorf("%s: %w", op, err)
	}

	return idx, nil
}
