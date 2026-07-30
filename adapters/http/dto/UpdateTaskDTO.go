package dto

import "time"

type UpdateTaskRequest struct {
	NAME        *string    `json:"name,omitempty"`
	DESCRIPTION *string    `json:"description,omitempty"`
	DEADLINE    *time.Time `json:"deadline,omitempty"`
}

type UpdateTaskResponse struct {
	STATUS string `json:"status"`
	ERROR  string `json:"error,omitempty"`
	ID     int64  `json:"id,omitempty"`
}
