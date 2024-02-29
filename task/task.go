package task

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Task represents a task to be done by a user.
type Task struct {
	ID        uuid.UUID
	Content   string
	CreatedAt time.Time
	DoneAt    time.Time
}
