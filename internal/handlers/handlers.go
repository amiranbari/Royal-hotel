package handlers

import (
	"github.com/amiranbari/Royal-hotel/internal/config"
	"github.com/amiranbari/Royal-hotel/internal/driver"
	"github.com/amiranbari/Royal-hotel/internal/models"
	"github.com/amiranbari/Royal-hotel/internal/renders"
	"github.com/amiranbari/Royal-hotel/internal/repository"
	"github.com/amiranbari/Royal-hotel/internal/repository/dbrepo"
	"net/http"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepo(app *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewPostgresDBRepo(db.SQL, app),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Index(rw http.ResponseWriter, r *http.Request) {
	renders.Template(rw, r, "index.page.tmpl", &models.TemplateData{})
}

func (m *Repository) SearchForRoom(rw http.ResponseWriter, r *http.Request) {
	
}
