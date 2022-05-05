package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/amiranbari/Royal-hotel/internal/config"
	"github.com/amiranbari/Royal-hotel/internal/renders"
	"log"
	"net/http"
	"time"
)

const portNumber string = ":8000"

var app config.AppConfig

func main() {
	app.InProduction = false

	/*
		| make session manager and put it into app config
	*/
	session := scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := renders.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	app.TemplateCache = tc
	app.UseCache = false

	fmt.Println(fmt.Sprintf("starting application on port number %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: route(&app),
	}

	err = srv.ListenAndServe()

	log.Fatal(err)
}
