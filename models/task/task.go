package task

import (
	"errors"
	"time"
)

// @name Task
type Task struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	DueDate     time.Time `json:"due_date" db:"due_date"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Validate checks if the task data is valid.
func (t *Task) Validate() error {
	// Check if the title is empty.
	if t.Title == "" || t.Description == "" {
		return errors.New("title or description dis required")
	}
	// TODO: дополонить валидацию

	return nil
}
