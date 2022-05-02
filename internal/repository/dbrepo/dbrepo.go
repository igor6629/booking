package dbrepo

import (
	"database/sql"
	"github.com/igor6629/booking/internal/config"
	"github.com/igor6629/booking/internal/repository"
)

type postgersDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgersDBRepo{
		App: a,
		DB:  conn,
	}
}
