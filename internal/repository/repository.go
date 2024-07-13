package repository

import (
	"hard/internal/domain/project"
	"hard/internal/domain/task"
	"hard/internal/domain/user"
	"hard/internal/repository/postgres"
	"hard/pkg/store"
)

type Configuration func(r *Repository) error

type Repository struct {
	postgres store.SQLX

	User    user.Repository
	Task    task.Repository
	Project project.Repository
}

func New(configs ...Configuration) (s *Repository, err error) {
	s = &Repository{}

	for _, cfg := range configs {
		if err = cfg(s); err != nil {
			return
		}
	}

	return
}

func WithPostgresStore(dbName string) Configuration {
	return func(r *Repository) (err error) {
		r.postgres, err = store.New(dbName)
		if err != nil {
			return
		}
		if err = store.Migrate(dbName); err != nil {
			return
		}

		r.User = postgres.NewUserRepository(r.postgres.Client)
		r.Task = postgres.NewTaskRepository(r.postgres.Client)
		r.Project = postgres.NewProjectRepository(r.postgres.Client)
		return
	}
}
