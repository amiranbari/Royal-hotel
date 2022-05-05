package main

import (
	"github.com/amiranbari/Royal-hotel/internal/config"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func route(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	return mux
}
