package store

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"strings"
)

func Migrate(dataSourceName string) (err error) {
	if !strings.Contains(dataSourceName, "://") {
		err = errors.New("store: undefined data source name " + dataSourceName)
		return
	}
	//driverName := strings.ToLower(strings.Split(dataSourceName, "://")[0])

	migrations, err := migrate.New("file://db/migrations", dataSourceName)
	if err != nil {
		return
	}

	if err = migrations.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
	}

	return
}

//psql -h localhost -U postgres -w -c "create database todo;"
//export POSTGRESQL_URL='postgres://postgres:postgres@localhost:5432/todo?sslmode=disable'
//export POSTGRESQL_URL='postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable'
//migrate create -ext sql -dir db/migrations -seq create_tasks_table
//sudo -u postgres createdb todo
