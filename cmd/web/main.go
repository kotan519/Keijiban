package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/kotan519/keijiban/internal/config"
	"github.com/kotan519/keijiban/internal/driver"
	"github.com/kotan519/keijiban/internal/handlers"
	"github.com/kotan519/keijiban/internal/helpers"
	"github.com/kotan519/keijiban/internal/models"
	"github.com/kotan519/keijiban/internal/render"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

// main is the main application function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.User{})
	gob.Register(models.TokumeiPostData{})
	gob.Register(models.TokumeiPostDataNumber{})
	gob.Register(models.TokumeiPostComment{})

	session = scs.New()

	app.Session = session

	log.Println("Connectiong to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=chiebukuro user=kotan519 password=")
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}else{
		log.Println("Connected to database")
	}
	
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)
	helpers.NewHelpers(&app)


	return db, nil
}

