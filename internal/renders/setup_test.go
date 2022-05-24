package renders

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/amiranbari/Royal-hotel/internal/config"
	"github.com/amiranbari/Royal-hotel/internal/models"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {
	//Say what we need to put in out session
	gob.Register(models.Reservation{})

	testApp.InProduction = false

	/*
		| make session manager and put it into testApp config

	*/
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction
	testApp.Session = session

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog
	testApp.UseCache = false

	app = &testApp
	os.Exit(m.Run())
}

type myWriter struct{}

func (mw myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (mw myWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func (mw myWriter) WriteHeader(statusCode int) {

}
