package user

import (
	"context"
	"hard/internal/domain/task"
)

type Repository interface {
	List(ctx context.Context) (dest []Entity, err error)
	Add(ctx context.Context, data Entity) (id string, err error)
	Get(ctx context.Context, id string) (data Entity, err error)
	Update(ctx context.Context, id string, data Entity) (err error)
	Delete(ctx context.Context, id string) (err error)
	Search(ctx context.Context, name string, email string) (data []Entity, err error)
	ListTasks(ctx context.Context, id string) (data []task.Entity, err error)
}

/*
GET /users: получить список всех пользователей.
POST /users: создать нового пользователя.
GET /users/{id}: получить данные конкретного пользователя.
PUT /users/{id}: обновить данные конкретного пользователя.
DELETE /users/{id}: удалить конкретного пользователя.
GET /users/{id}/tasks: получить список задач конкретного пользователя.
GET /users/search?name={name}: найти пользователей по имени.
GET /users/search?email={email}: найти пользователей по электронной почте.
*/
