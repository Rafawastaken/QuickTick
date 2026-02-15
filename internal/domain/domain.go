package domain

import "time"

type Status string

const (
	StatusTodo      Status = "todo"
	StatusProgress  Status = "progress"
	StatusCompleted Status = "completed"
	StatusCanceled  Status = "canceled"
)

func IsValidStatus(s Status) bool {
	switch s {
	case StatusTodo, StatusProgress, StatusCompleted, StatusCanceled:
		return true
	default:
		return false
	}
}

type Task struct {
	ID        int64
	Title     string
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Filter struct {
	Status Status
}
