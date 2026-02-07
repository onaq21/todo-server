package task

import "time"

type Task struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}
