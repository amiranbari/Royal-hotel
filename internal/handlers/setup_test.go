package handlers

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/amiranbari/Royal-hotel/internal/config"
	"github.com/amiranbari/Royal-hotel/internal/helpers"
	"github.com/amiranbari/Royal-hotel/internal/models"
	"github.com/amiranbari/Royal-hotel/internal/renders"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger
var functions = template.FuncMap{}

func TestMain(m *testing.M) {
	//Say what we need to put in out session
	gob.Register(models.Reservation{})

	app.InProduction = false

	/*
		| make session manager and put it into app config

	*/
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	//Make template cache
	tc, err := createTestTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := NewTestRepo(&app)
	NewHandlers(repo)
	renders.NewRenderer(&app)
	helpers.NewHelpers(&app)

	os.Exit(m.Run())
}

func createTestTemplateCache() (config.TemplateCache, error) {
	myCache := config.TemplateCache{}

	pages, err := filepath.Glob("../../templates/*.page.tmpl")

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("../../templates/*.layout.tmpl")

		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("../../templates/*.layout.tmpl")

			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}

func getRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Index)
	mux.Get("/search", Repo.SearchForRoom)
	mux.Get("/login", Repo.Login)
	mux.Post("/login", Repo.PostLogin)
	mux.Get("/logout", Repo.Logout)

	mux.Route("/", func(mux chi.Router) {
		mux.Use(UserAuth)
		mux.Get("/choose-room/{id}", Repo.ChooseRoom)
		mux.Get("/make-reservation", Repo.MakeReservation)
		mux.Post("/make-reservation", Repo.PostReservation)
	})

	//static files
	fileServer := http.FileServer(http.Dir("../../static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}

func NoServe(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

//SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

//UserAuth router only authenticated users
func UserAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if !helpers.IsUserAuthenticated(r) {
			session.Put(r.Context(), "error", "Log in first!")
			http.Redirect(rw, r, "/login?redirect="+r.URL.Path, http.StatusSeeOther)
			return
		}
		next.ServeHTTP(rw, r)
	})
}
