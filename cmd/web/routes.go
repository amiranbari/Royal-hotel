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
	mux.Get("/search", handlers.Repo.SearchForRoom)
	mux.Get("/login", handlers.Repo.Login)
	mux.Post("/login", handlers.Repo.PostLogin)
	mux.Get("/logout", handlers.Repo.Logout)

	mux.Route("/", func(mux chi.Router) {
		mux.Use(UserAuth)
		mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
		mux.Get("/make-reservation", handlers.Repo.MakeReservation)
		mux.Post("/make-reservation", handlers.Repo.PostReservation)
	})

	//static files
	fileServer := http.FileServer(http.Dir("../../static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
