package main

import (
	"github.com/amiranbari/Royal-hotel/internal/config"
	"github.com/amiranbari/Royal-hotel/internal/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func route(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", handlers.Repo.Index)

	return mux
}
