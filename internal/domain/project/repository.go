package project

import (
	"context"
	"hard/internal/domain/task"
)

type Repository interface {
	List(ctx context.Context) (dest []Entity, err error)
	Add(ctx context.Context, data Entity) (id string, err error)
	Get(ctx context.Context, id string) (dest Entity, err error)
	Update(ctx context.Context, id string, dest Entity) (err error)
	Delete(ctx context.Context, id string) (err error)
	Search(ctx context.Context, data Entity) (dest []Entity, err error)
	ListTasks(ctx context.Context, id string) (data []task.Entity, err error)
}
