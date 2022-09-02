package render

import (
	"fmt"
	"html/template"
	"net/http"
)

// Renders template to ResponseWrite using file extension
func RenderTemplate(w http.ResponseWriter, tmpl string) {
	// Use parse files to grab template files
	parsedTemplate, parseError := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.tmpl")
	// If error detected
	if parseError != nil {
		fmt.Println("Error encountered with file parse:", parseError.Error())
		return
	}
	// Execute the template by injecting data
	err := parsedTemplate.Execute(w, nil)

	// If error detected
	if err != nil {
		fmt.Println("Error encountered with render:")
		return
	}
}

var templateCache = make(map[string]*template.Template)

// Checks if template has already been rendered in the cache
func RenderTemplateTest(w http.ResponseWriter, t string) {
	var templateToRender *template.Template
	var err error

	// check to see if we have the template in the cache
	// if string not present in cache, it will return false
	_, inMap := templateCache[t]

	// If not found, create
	if !inMap {

		createTemplateCacheEntry(t)

	}

	// Grab template from map and execute render
	templateToRender = templateCache[t]
	renderErr := templateToRender.Execute(w, nil)

	// If error detected
	if renderErr != nil {
		fmt.Println("Error encountered with newly found template execution:", err.Error())
		return
	}
}

// Creates template cache using requested file and base layout
func createTemplateCacheEntry(t string) error {
	// Create string slice containing base layout and requested layout
	templates := []string{

		fmt.Sprintf("./templates/%s", t),
		"./templates/base.layout.tmpl",
	}

	// Unpack string slice in parsefiles to grab all templates
	tmpl, err := template.ParseFiles(templates...)

	// If error found
	if err != nil {
		return err
	}
	// Save parsed template in template cache
	fmt.Println("adding parsed template to cache")
	templateCache[t] = tmpl

	return nil
}
