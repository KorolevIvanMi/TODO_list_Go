package dto

import "time"

type GetAllTasksRespModel struct {
	ID          int
	NAME        string
	DESCRIPTION string
	DEADLINE    time.Time
}

type GetAllTasksResponse struct {
	STATUS string                 `json:"status"`
	TASKS  []GetAllTasksRespModel `json:"tasks,omitempty"`
	AMOUNT uint64                 `json:"amount,omitempty"`
}
