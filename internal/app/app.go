package app

import (
	"context"

	"fmt"
	"os"

	"github.com/rafawastaken/quicktick/internal/config"
	"github.com/rafawastaken/quicktick/internal/domain"
	"github.com/rafawastaken/quicktick/internal/notes"
	"github.com/rafawastaken/quicktick/internal/store"
	"github.com/rafawastaken/quicktick/internal/store/sqlite"
	"github.com/rafawastaken/quicktick/internal/sync"
	"github.com/rafawastaken/quicktick/internal/sync/supabase"
)

type App struct {
	Store  store.Store
	Syncer sync.Syncer
}

func New(ctx context.Context, userID string) (*App, error) {
	dbPath, err := config.DBPath(userID)
	if err != nil {
		return nil, err
	}

	st, err := sqlite.Open(dbPath)
	if err != nil {
		return nil, err
	}

	if err := st.Init(ctx); err != nil {
		st.Close()
		return nil, err
	}

	// For now, use NoOpSyncer
	// syncer := sync.NewNoOpSyncer()

	// Use Supabase Syncer
	syncer := supabase.NewSyncer(st) // We need to export NewSyncer or similar

	// BUT wait, circular dependency?
	// app -> supabase -> store (good)
	// app -> sync (good)
	// We need to import internal/sync/supabase in app/app.go

	return &App{
		Store:  st,
		Syncer: syncer,
	}, nil
}

func (a *App) Close() error {
	return a.Store.Close()
}

func (a *App) AddTask(ctx context.Context, title string, status domain.Status) (int64, error) {
	return a.Store.AddTask(ctx, title, status)
}

func (a *App) ListTasks(ctx context.Context, filter domain.Filter) ([]domain.Task, error) {
	return a.Store.ListTasks(ctx, filter)
}

func (a *App) CompleteTask(ctx context.Context, id int64) error {
	return a.Store.UpdateStatus(ctx, id, domain.StatusCompleted)
}

func (a *App) DeleteTask(ctx context.Context, id int64) error {
	return a.Store.DeleteTask(ctx, id)
}

func (a *App) EditTask(ctx context.Context, id int64, title string) error {
	return a.Store.EditTaskTitle(ctx, id, title)
}

func (a *App) SyncTasks(ctx context.Context) error {
	return a.Syncer.Sync(ctx)
}

func (a *App) OpenTask(ctx context.Context, id int64) error {
	t, err := a.Store.GetTask(ctx, id)
	if err != nil {
		return err
	}

	notePath, err := notes.GetNotePath(id)
	if err != nil {
		return err
	}

	if _, err := os.Stat(notePath); os.IsNotExist(err) {
		content := fmt.Sprintf("# Note for Task #%d: %s\n\n", t.ID, t.Title)
		if err := os.WriteFile(notePath, []byte(content), 0644); err != nil {
			return err
		}
	}

	return notes.OpenEditor(notePath)
}
