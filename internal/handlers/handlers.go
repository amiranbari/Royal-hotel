package handlers

import (
	"github.com/amiranbari/Royal-hotel/internal/config"
	"net/http"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

func NewRepo(app *config.AppConfig) *Repository {
	return &Repository{
		App: app,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Index(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("welcome"))
}
