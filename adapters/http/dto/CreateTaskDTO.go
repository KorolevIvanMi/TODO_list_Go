package dto

import "time"

type CreateTaskRequst struct {
	NAME        string    `json:"name"`
	DESCRIPTION string    `json:"description"`
	DEADLINE    time.Time `json:"deadline"`
}

type CreateTaskResponse struct {
	STATUS string `json:"status"`
	ERROR  string `json:"error,omitempty"`
	NAME   string `json:"name,omitempty"`
	ID     int64  `json:"id,omitempty"`
}
