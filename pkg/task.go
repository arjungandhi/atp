package atp

import (
	"time"

	"github.com/google/uuid"
)

// single task type
type Task struct {
	Id                uuid.UUID     `json:"id"`
	Name              string        `json:"name"`
	Description       string        `json:"description"`
	Duration          time.Duration `json:"duration"`
	EstimatedDuration time.Duration `json:"estimated_duration"`
	Deadline          time.Time     `json:"deadline"`
	Completed         bool          `json:"completed"`
}

// Get a new Task
func NewTask() {
}
