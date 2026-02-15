package supabase

import (
	"time"

	"github.com/rafawastaken/quicktick/internal/domain"
)

type TaskJSON struct {
	ID        int64     `json:"id,omitempty"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    string    `json:"user_id,omitempty"` // For RLS we might need to send this, or not?
}

func ToJSON(t domain.Task) TaskJSON {
	return TaskJSON{
		ID:        t.ID,
		Title:     t.Title,
		Status:    string(t.Status),
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func ToDomain(j TaskJSON) domain.Task {
	return domain.Task{
		ID:        j.ID,
		Title:     j.Title,
		Status:    domain.Status(j.Status),
		CreatedAt: j.CreatedAt,
		UpdatedAt: j.UpdatedAt,
	}
}
