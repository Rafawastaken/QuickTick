package store

import (
	"context"

	"github.com/rafawastaken/quicktick/internal/domain"
)

type Store interface {
	Init(ctx context.Context) error

	AddTask(ctx context.Context, title string, status domain.Status) (int64, error)
	ListTasks(ctx context.Context, filter domain.Filter) ([]domain.Task, error)
	UpdateStatus(ctx context.Context, id int64, status domain.Status) error
	EditTaskTitle(ctx context.Context, id int64, title string) error
	DeleteTask(ctx context.Context, id int64) error
	GetTask(ctx context.Context, id int64) (domain.Task, error)
	Close() error
}
