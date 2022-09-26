package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dmawardi/bookings/internal/config"
	"github.com/dmawardi/bookings/internal/handlers"
	"github.com/dmawardi/bookings/internal/helpers"
	"github.com/dmawardi/bookings/internal/models"
	"github.com/dmawardi/bookings/internal/render"
)

// Init state (incl. templates)
var app config.AppConfig

const portNumber = ":8080"

var session *scs.SessionManager

// Loggers
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Starting application on port: %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {
	// In Session, you can store primitives.
	// If more complex, define first
	gob.Register(models.ReservationForm{})

	// Change this to true in production
	app.InProduction = false

	// Create info log that outputs to std output, has prefix INFO, then date/time
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// Set app info log as created logger
	app.InfoLog = infoLog

	// Create error log that outputs to std output, has prefix INFO, then date/time
	// short file provides information on error
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	// Set app error log as created logger
	app.ErrorLog = errorLog

	session = scs.New()
	// Set session lifetime to 24 hours
	session.Lifetime = 24 * time.Hour
	// Persist when browser restarted
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	// For https set to true
	session.Cookie.Secure = app.InProduction
	// Set session to state variable: Session
	app.Session = session

	// create cache
	createdCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Failed to create template cache")
		return err
	}
	// Store created cache in app config
	app.TemplateCache = createdCache
	// Config: Use cache
	// if false, will serve updated file (dev mode)
	app.UseCache = false

	// Create new handler repository
	repo := handlers.NewRepo(&app)
	// Build new handlers
	handlers.UpdateRepositoryHandlers(repo)

	// Sets template cache for render package
	render.SetTemplate(&app)
	helpers.SetState(&app)
	return nil
}
