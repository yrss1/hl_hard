package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"hard/internal/domain/task"
	"hard/pkg/store"
	"strings"
)

type TaskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) List(ctx context.Context) (dest []task.Entity, err error) {
	query := `
			SELECT id, title, description, priority, status, assignee_id, project_id, completed_at
			FROM tasks
			ORDER BY id`

	err = r.db.SelectContext(ctx, &dest, query)

	return
}

func (r *TaskRepository) Add(ctx context.Context, data task.Entity) (id string, err error) {
	query := `
		INSERT INTO tasks (title, description, priority, status, assignee_id, project_id, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING id`

	args := []any{data.Title, data.Description, data.Priority, data.Status, data.AssigneeID, data.ProjectID, data.CompletedAt}

	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *TaskRepository) Get(ctx context.Context, id string) (dest task.Entity, err error) {
	query := `
		SELECT id, title, description, priority, status, assignee_id, project_id, completed_at 
		FROM tasks 
		WHERE id=$1`

	args := []any{id}

	if err = r.db.GetContext(ctx, &dest, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *TaskRepository) Update(ctx context.Context, id string, data task.Entity) (err error) {
	sets, args := r.prepareArgs(data)
	if len(args) > 0 {
		args = append(args, id)
		sets = append(sets, "updated_at=CURRENT_TIMESTAMP")

		query := fmt.Sprintf("UPDATE tasks SET %s WHERE id=$%d RETURNING id", strings.Join(sets, ", "), len(args))
		if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = store.ErrorNotFound
			}
		}
	}

	return
}

func (r *TaskRepository) prepareArgs(data task.Entity) (sets []string, args []any) {
	if data.Title != nil {
		args = append(args, data.Title)
		sets = append(sets, fmt.Sprintf("title=$%d", len(args)))
	}

	if data.Description != nil {
		args = append(args, data.Description)
		sets = append(sets, fmt.Sprintf("description=$%d", len(args)))
	}

	if data.Priority != nil {
		args = append(args, data.Priority)
		sets = append(sets, fmt.Sprintf("priority=$%d", len(args)))
	}

	if data.Status != nil {
		args = append(args, data.Status)
		sets = append(sets, fmt.Sprintf("status=$%d", len(args)))
	}

	if data.AssigneeID != nil {
		args = append(args, data.AssigneeID)
		sets = append(sets, fmt.Sprintf("assignee_id=$%d", len(args)))
	}

	if data.ProjectID != nil {
		args = append(args, data.ProjectID)
		sets = append(sets, fmt.Sprintf("project_id=$%d", len(args)))
	}

	if data.CompletedAt != nil {
		args = append(args, *data.CompletedAt)
		sets = append(sets, fmt.Sprintf("completed_at=$%d", len(args)))
	}

	return
}

func (r *TaskRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
		DELETE FROM tasks
		WHERE id=$1
		RETURNING id`

	args := []any{id}

	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *TaskRepository) Search(ctx context.Context, data task.Entity) (dest []task.Entity, err error) {
	query := "SELECT id, title, description, priority, status, assignee_id, project_id, completed_at FROM tasks WHERE 1=1"

	sets, args := r.prepareArgs(data)
	if len(sets) > 0 {
		query += " AND " + strings.Join(sets, " AND ")
	}

	err = r.db.SelectContext(ctx, &dest, query, args...)

	return
}
