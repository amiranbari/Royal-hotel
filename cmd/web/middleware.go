package main

import (
	"github.com/amiranbari/Royal-hotel/internal/helpers"
	"github.com/justinas/nosurf"
	"net/http"
)

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
