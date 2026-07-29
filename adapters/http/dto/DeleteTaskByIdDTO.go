package dto

type DeleteTaskByIdResponse struct {
	STATUS string `json:"status"`
	ID     int    `json:"id,omitempty"`
}
