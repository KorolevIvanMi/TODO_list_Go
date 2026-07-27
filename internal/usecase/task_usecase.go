package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/KorolevIvanMi/TODO_list_Go/internal/domain"
	"github.com/KorolevIvanMi/TODO_list_Go/internal/repository"
)

type UseCase struct {
	repo repository.TaskRepository
}

type CreatedTask struct {
	ID          int64
	Name        string
	Description string
	Deadline    time.Time
}

func New(repo repository.TaskRepository) *UseCase {
	var uc UseCase
	uc.repo = repo
	return &uc
}

func (uc *UseCase) CreateTask(ctx context.Context, name, description string, deadline time.Time) (*CreatedTask, error) {

	const op = "internal.usecase.task.usecase.CreateTask"
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	task, err := domain.NewTask(name, description, deadline)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	idx, err := uc.repo.SaveTask(ctx, task.Name, task.Description, task.Deadline)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	newTask := CreatedTask{ID: idx, Name: task.Name, Description: task.Description, Deadline: task.Deadline}
	return &newTask, nil

}
