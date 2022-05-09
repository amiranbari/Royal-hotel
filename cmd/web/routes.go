package main

import (
	"github.com/amiranbari/Royal-hotel/internal/config"
	"github.com/amiranbari/Royal-hotel/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func route(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoServe)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Index)
	mux.Post("/", handlers.Repo.SearchForRoom)

	//static files
	fileServer := http.FileServer(http.Dir("../../static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
