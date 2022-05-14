package helpers

import (
	"fmt"
	"github.com/amiranbari/Royal-hotel/internal/config"
	"net/http"
	"runtime/debug"
)

var app *config.AppConfig

// NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(rw http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(rw, http.StatusText(status), status)
}

func ServerError(rw http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n", debug.Stack())
	app.ErrorLog.Println(err)
	app.ErrorLog.Println(trace)
	http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func IsUserAuthenticated(r *http.Request) bool {
	return app.Session.Exists(r.Context(), "user_id")
}
