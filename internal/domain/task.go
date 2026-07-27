package domain

import (
	"fmt"
	"time"
)

type Task struct {
	Name        string
	Description string
	Deadline    time.Time
}

func NewTask(name, description string, deadline time.Time) (*Task, error) {
	const op = "internal.domain.NewTask"

	if name == "" || len(name) > 100 {
		return nil, fmt.Errorf("%s : not correct name", op)
	}

	if description == "" {
		return nil, fmt.Errorf("%s : not correct description", op)
	}

	if deadline.Before(time.Now()) {
		return nil, fmt.Errorf("%s : deadline can not be in the past", op)
	}

	new_task := Task{Name: name, Description: description, Deadline: deadline}

	return &new_task, nil
}
