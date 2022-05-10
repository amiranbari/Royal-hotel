package handlers

import (
	"github.com/amiranbari/Royal-hotel/internal/config"
	"github.com/amiranbari/Royal-hotel/internal/driver"
	"github.com/amiranbari/Royal-hotel/internal/forms"
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
	renders.Template(rw, r, "index.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) SearchForRoom(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "form is not valid!")
		http.Redirect(rw, r, "/", http.StatusSeeOther)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("start_date", "end_date")
	if !form.Valid() {
		renders.Template(rw, r, "index.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

}
