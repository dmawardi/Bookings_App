package handlers

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"html/template"

	"github.com/alexedwards/scs/v2"
	"github.com/dmawardi/bookings/internal/config"
	"github.com/dmawardi/bookings/internal/models"
	"github.com/dmawardi/bookings/internal/render"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager

// Loggers
var infoLog *log.Logger
var errorLog *log.Logger

func getRoutes() http.Handler {
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
	createdCache, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("Failed to create template cache")
	}
	// Store created cache in app config
	app.TemplateCache = createdCache
	// Config: Use cache
	// if false, will serve updated file (dev mode)
	app.UseCache = true

	// Create new handler repository
	repo := NewRepo(&app)
	// Build new handlers
	UpdateRepositoryHandlers(repo)

	// Sets template cache for render package
	render.SetTemplate(&app)

	// Setup middleware and routes
	//
	// Create new Chi router
	mux := chi.NewRouter()

	// Use built in Chi middleware
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)
	// Ignore any POST request that doesn't have CSRF token
	// mux.Use(NoSurfCSRF)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	mux.Get("/majors-suite", Repo.Majors)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/contact", Repo.Contact)

	// Build fileserver using static directory
	fileServer := http.FileServer(http.Dir("./static"))
	// Handle all calls to /static/* by stripping prefix and sending to file server
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))
	// Return complete router
	return mux
}

// Middleware functions copied over as we can't import main package
// Adds CSRF protection for all POST requests
func NoSurfCSRF(next http.Handler) http.Handler {
	// Create a handler using nosurf and the next param
	csrfHandler := nosurf.New(next)

	// Set cookie settings
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		// Apply to entire site
		Path: "/",
		// Using https?
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// Loads and saves the session for every request
func SessionLoad(next http.Handler) http.Handler {
	// Load and save session for use
	return session.LoadAndSave(next)
}

var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	// templateCache := make(map[string]*template.Template)
	templateCache := map[string]*template.Template{}

	// get all files with page.tmpl from templates folder
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		// Return the cache and error
		return templateCache, err
	}

	// range through all pages found
	for _, pagePath := range pages {
		// Base returns the last element of path (ie. filename)
		fileName := filepath.Base(pagePath)

		// name template as filename and parsefile
		templateSet, err := template.New(fileName).Funcs(functions).ParseFiles(pagePath)
		if err != nil {
			fmt.Println("error encountered building template set.", err.Error())
			return templateCache, err
		}

		// get all files with layout.tmpl from templates folder
		layoutMatches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			// Return the cache and error
			return templateCache, err
		}

		// if any layoutMatches are found
		if len(layoutMatches) > 0 {
			// Adds layoutMatches to template set using parseGlob
			templateSet, err = templateSet.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				// Return the cache and error
				return templateCache, err
			}
		}
		// Add template set to myCache
		templateCache[fileName] = templateSet
	}

	return templateCache, nil
}
