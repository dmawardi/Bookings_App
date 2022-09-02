package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/dmawardi/bookings/internal/config"
	"github.com/dmawardi/bookings/internal/models"
	"github.com/justinas/nosurf"
)

// Used when parsing files
var functions = template.FuncMap{}

var app *config.AppConfig

// Sets the config for the template page
func SetTemplate(a *config.AppConfig) {
	app = a
}

// Adds default data for every page
func AddDefaultTemplateData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// Add defaults
	// td.StringMap["Sample"] = "Sample default data."
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")

	// Set CSRF token
	td.CSRFToken = nosurf.Token(r)
	return td
}

func AltRenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	// Create template cache var
	var templateCache map[string]*template.Template
	// If config detected to use cache
	if app.UseCache {

		// Grab cache from app config
		templateCache = app.TemplateCache
	} else {
		// rebuild the template cache
		templateCache, _ = CreateTemplateCache()
	}

	// Find template in cache
	foundTemplate, templateError := templateCache[tmpl]
	if !templateError {
		log.Fatal("Could not retrieve template from cache", templateError)
	}

	// Create new buffer
	buf := new(bytes.Buffer)

	AddDefaultTemplateData(td, r)
	// Execute template in buffer using data
	dataInputError := foundTemplate.Execute(buf, td)
	if dataInputError != nil {
		fmt.Println("data input error", dataInputError.Error())
	}
	// Write above execution to response writer
	_, executionError := buf.WriteTo(w)

	// If error detected
	if executionError != nil {
		fmt.Println("Error encountered with render.", executionError.Error())
		return
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	// templateCache := make(map[string]*template.Template)
	templateCache := map[string]*template.Template{}

	// get all files with page.tmpl from templates folder
	pages, err := filepath.Glob("./templates/*.page.tmpl")
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
		layoutMatches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			// Return the cache and error
			return templateCache, err
		}

		// if any layoutMatches are found
		if len(layoutMatches) > 0 {
			// Adds layoutMatches to template set using parseGlob
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.tmpl")
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
