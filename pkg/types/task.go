package types

import "time"

type Task struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"updated_at"`
	Description string    `json:"description"`
	ID          int64     `json:"id"`
	Done        bool      `json:"done"`
}
