package repository

import (
	"context"
	"time"
)

type TaskRepository interface {
	SaveTask(ctx context.Context, name, description string, deadline time.Time) (int64, error)
}
