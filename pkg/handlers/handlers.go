package handlers

import (
	"net/http"

	"github.com/dmawardi/bookings/pkg/config"
	"github.com/dmawardi/bookings/pkg/models"
	"github.com/dmawardi/bookings/pkg/render"
)

// Repository used by handler package
var Repo *Repository

// Repository type
type Repository struct {
	App *config.AppConfig
}

// Create new handler repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// Set Repository to parameter
func UpdateRepositoryHandlers(r *Repository) {
	Repo = r
}

// Home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// Store IP of user
	remoteIP := r.RemoteAddr
	// Add a key to session to data
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	stringMap := make(map[string]string)
	stringMap["boo"] = "boodoo"

	render.AltRenderTemplate(w, "home.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// About page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello again"
	stringMap["boo"] = "boodoo"
	// Get remote ip from session data
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.AltRenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})

}
