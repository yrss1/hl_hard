package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"hard/internal/domain/task"
	"hard/internal/domain/user"
	"hard/pkg/store"
	"strings"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}
func (r *UserRepository) List(ctx context.Context) (dest []user.Entity, err error) {
	query := `
			SELECT id, full_name, email, role 
			FROM users
			ORDER BY id`

	err = r.db.SelectContext(ctx, &dest, query)

	return
}

func (r *UserRepository) Add(ctx context.Context, data user.Entity) (id string, err error) {
	query := `
		INSERT INTO users (full_name, email, role) 
		VALUES ($1, $2, $3) 
		RETURNING id`
	args := []any{data.FullName, data.Email, data.Role}
	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}
	return
}

func (r *UserRepository) Get(ctx context.Context, id string) (dest user.Entity, err error) {
	query := `
		SELECT id, full_name, email, role
		FROM users 
		WHERE id=$1`

	args := []any{id}

	if err = r.db.GetContext(ctx, &dest, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}
	return
}

func (r *UserRepository) Update(ctx context.Context, id string, data user.Entity) (err error) {
	sets, args := r.prepareArgs(data)
	if len(args) > 0 {

		args = append(args, id)
		sets = append(sets, "updated_at=CURRENT_TIMESTAMP")
		query := fmt.Sprintf("UPDATE users SET %s WHERE id=$%d RETURNING id", strings.Join(sets, ", "), len(args))

		if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = store.ErrorNotFound
			}
		}
	}

	return
}

func (r *UserRepository) prepareArgs(data user.Entity) (sets []string, args []any) {
	if data.FullName != nil {
		args = append(args, *data.FullName)
		sets = append(sets, fmt.Sprintf("full_name=$%d", len(args)))
	}

	if data.Email != nil {
		args = append(args, *data.Email)
		sets = append(sets, fmt.Sprintf("email=$%d", len(args)))
	}

	if data.Role != nil {
		args = append(args, *data.Role)
		sets = append(sets, fmt.Sprintf("role=$%d", len(args)))
	}
	return
}

func (r *UserRepository) Delete(ctx context.Context, id string) (err error) {
	args := []any{id}

	updateQuery := `
        UPDATE tasks
        SET assignee_id = NULL
        WHERE assignee_id = $1
    `
	if err = r.db.QueryRowContext(ctx, updateQuery, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	deleteQuery := `
        DELETE FROM users
        WHERE id = $1
    `
	if err = r.db.QueryRowContext(ctx, deleteQuery, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *UserRepository) ListTasks(ctx context.Context, id string) (dest []task.Entity, err error) {
	args := []any{id}
	existsQuery := `
			SELECT 1
			FROM users 
			WHERE id=$1
		`

	if err = r.db.QueryRowContext(ctx, existsQuery, id).Scan(&id); err != nil {
		err = store.ErrorNotFound
		return
	}

	query := `
		SELECT id, title, description, priority, status, assignee_id, project_id, completed_at
		FROM tasks 
		WHERE assignee_id=$1`

	if err = r.db.SelectContext(ctx, &dest, query, args...); err != nil {
		return
	}

	return
}

func (r *UserRepository) Search(ctx context.Context, name string, email string) (dest []user.Entity, err error) {
	sets, args := r.prepareSearchArgs(name, email)
	query := fmt.Sprintf("SELECT id, full_name, email, role FROM users WHERE 1=1 %s", strings.Join(sets, " "))

	err = r.db.SelectContext(ctx, &dest, query, args...)

	return dest, nil
}

func (r *UserRepository) prepareSearchArgs(name string, email string) (sets []string, args []any) {
	if name != "" {
		args = append(args, "%"+name+"%")
		sets = append(sets, fmt.Sprintf("AND full_name ILIKE $%d", len(args)))
	}
	if email != "" {
		args = append(args, "%"+email+"%")
		sets = append(sets, fmt.Sprintf("AND email ILIKE $%d", len(args)))
	}
	return
}
