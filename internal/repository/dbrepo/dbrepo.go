package dbrepo

import (
	"database/sql"
	"github.com/amiranbari/Royal-hotel/internal/config"
	"github.com/amiranbari/Royal-hotel/internal/repository"
)

type PostgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresDBRepo(conn *sql.DB, app *config.AppConfig) repository.DatabaseRepo {
	return &PostgresDBRepo{
		App: app,
		DB:  conn,
	}
}

func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
