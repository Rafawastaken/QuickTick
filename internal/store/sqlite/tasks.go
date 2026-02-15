package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/rafawastaken/quicktick/internal/domain"
	"github.com/rafawastaken/quicktick/internal/util"
)

func (s *Store) AddTask(ctx context.Context, title string, status domain.Status) (int64, error) {
	now := util.FormatTime(util.Now())

	res, err := s.db.ExecContext(ctx,
		`INSERT INTO tasks(title, status, created_at, updated_at) VALUES(?, ?, ?, ?)`,
		title, string(status), now, now,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) ListTasks(ctx context.Context, filter domain.Filter) ([]domain.Task, error) {
	query := `SELECT id, title, status, created_at, updated_at FROM tasks`
	var args []interface{}

	if filter.Status != "" {
		query += ` WHERE status = ?`
		args = append(args, string(filter.Status))
	}

	query += ` ORDER BY id DESC`

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Task
	for rows.Next() {
		var (
			t               domain.Task
			status          string
			createdAtString string
			updatedAtString string
		)
		if err := rows.Scan(&t.ID, &t.Title, &status, &createdAtString, &updatedAtString); err != nil {
			return nil, err
		}
		t.Status = domain.Status(status)
		t.CreatedAt, _ = util.ParseTime(createdAtString)
		t.UpdatedAt, _ = util.ParseTime(updatedAtString)
		out = append(out, t)
	}
	return out, rows.Err()
}

func (s *Store) UpdateStatus(ctx context.Context, id int64, status domain.Status) error {
	now := util.FormatTime(util.Now())

	res, err := s.db.ExecContext(ctx,
		`UPDATE tasks SET status = ?, updated_at = ? WHERE id = ?`,
		string(status), now, id,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return util.ErrNotFound
	}
	return nil
}

func (s *Store) EditTaskTitle(ctx context.Context, id int64, title string) error {
	now := util.FormatTime(util.Now())
	res, err := s.db.ExecContext(ctx,
		`UPDATE tasks SET title = ?, updated_at = ? WHERE id = ?`,
		title, now, id,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return util.ErrNotFound
	}
	return nil
}

func (s *Store) DeleteTask(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return util.ErrNotFound
	}
	return nil
}

func (s *Store) GetTask(ctx context.Context, id int64) (domain.Task, error) {
	var (
		t               domain.Task
		status          string
		createdAtString string
		updatedAtString string
	)
	err := s.db.QueryRowContext(ctx,
		`SELECT id, title, status, created_at, updated_at FROM tasks WHERE id = ?`,
		id,
	).Scan(&t.ID, &t.Title, &status, &createdAtString, &updatedAtString)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Task{}, util.ErrNotFound
		}
		return domain.Task{}, err
	}

	t.Status = domain.Status(status)
	t.CreatedAt, _ = util.ParseTime(createdAtString)
	t.UpdatedAt, _ = util.ParseTime(updatedAtString)
	return t, nil
}
