package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/amiranbari/Royal-hotel/internal/config"
	"github.com/amiranbari/Royal-hotel/internal/driver"
	"github.com/amiranbari/Royal-hotel/internal/handlers"
	"github.com/amiranbari/Royal-hotel/internal/helpers"
	"github.com/amiranbari/Royal-hotel/internal/models"
	"github.com/amiranbari/Royal-hotel/internal/renders"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber string = ":8000"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("starting application on port number %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: route(&app),
	}

	err = srv.ListenAndServe()

	log.Fatal(err)
}

func run() (*driver.DB, error) {
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

	//connect to database
	log.Println("Connecting to database ...")
	db, err := driver.ConnectSql("host=localhost port=5432 dbname=royal-hotel user=postgres password=123456")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying ...", err)
	}
	log.Println("Connected to database!")

	//Make template cache
	tc, err := renders.CreateTemplateCache()
	if err != nil {
		return db, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	renders.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
