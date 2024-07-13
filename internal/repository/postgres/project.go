package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"hard/internal/domain/project"
	"hard/internal/domain/task"
	"hard/pkg/store"
	"strings"
)

type ProjectRepository struct {
	db *sqlx.DB
}

func NewProjectRepository(db *sqlx.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) List(ctx context.Context) (dest []project.Entity, err error) {
	query := `
			SELECT id, title, description, start_date, end_date, manager_id 
			FROM projects
			ORDER BY id`

	err = r.db.SelectContext(ctx, &dest, query)

	return
}

func (r *ProjectRepository) Add(ctx context.Context, data project.Entity) (id string, err error) {
	query := `
		INSERT INTO projects (title, description, start_date, end_date, manager_id ) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id`
	args := []any{data.Title, data.Description, data.StartDate, data.EndDate, data.ManagerID}
	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}
	return
}

func (r *ProjectRepository) Get(ctx context.Context, id string) (dest project.Entity, err error) {
	query := `
		SELECT id, title, description, start_date, end_date, manager_id
		FROM projects 
		WHERE id=$1`

	args := []any{id}

	if err = r.db.GetContext(ctx, &dest, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}
	return
}

func (r *ProjectRepository) Update(ctx context.Context, id string, data project.Entity) (err error) {
	sets, args := r.prepareArgs(data)
	if len(args) > 0 {

		args = append(args, id)
		sets = append(sets, "updated_at=CURRENT_TIMESTAMP")
		query := fmt.Sprintf("UPDATE projects SET %s WHERE id=$%d RETURNING id", strings.Join(sets, ", "), len(args))
		fmt.Println(query)
		if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = store.ErrorNotFound
			}
		}
	}

	return
}

func (r *ProjectRepository) prepareArgs(data project.Entity) (sets []string, args []any) {
	if data.Title != nil {
		args = append(args, data.Title)
		sets = append(sets, fmt.Sprintf("title=$%d", len(args)))
	}

	if data.Description != nil {
		args = append(args, data.Description)
		sets = append(sets, fmt.Sprintf("description=$%d", len(args)))
	}

	if !data.StartDate.IsZero() {
		args = append(args, data.StartDate)
		sets = append(sets, fmt.Sprintf("start_date=$%d", len(args)))
	}

	if !data.EndDate.IsZero() {
		args = append(args, data.EndDate)
		sets = append(sets, fmt.Sprintf("end_date=$%d", len(args)))
	}

	if data.ManagerID != nil {
		args = append(args, data.ManagerID)
		sets = append(sets, fmt.Sprintf("manager_id=$%d", len(args)))
	}

	return
}

func (r *ProjectRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
		DELETE FROM projects
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

func (r *ProjectRepository) Search(ctx context.Context, data project.Entity) (dest []project.Entity, err error) {
	query := "SELECT id, title, description, start_date, end_date, manager_id FROM projects WHERE 1=1"

	sets, args := r.prepareArgs(data)

	if len(sets) > 0 {
		query += " AND " + strings.Join(sets, " AND ")
	}
	err = r.db.SelectContext(ctx, &dest, query, args...)

	return
}

func (r *ProjectRepository) ListTasks(ctx context.Context, id string) (dest []task.Entity, err error) {
	args := []any{id}
	existsQuery := `
			SELECT 1
			FROM projects 
			WHERE id=$1
		`

	if err = r.db.QueryRowContext(ctx, existsQuery, id).Scan(&id); err != nil {
		err = store.ErrorNotFound
		return
	}

	query := `
		SELECT id, title, description, priority, status, assignee_id, project_id, completed_at
		FROM tasks 
		WHERE project_id=$1`

	if err = r.db.SelectContext(ctx, &dest, query, args...); err != nil {
		return
	}

	return
}
