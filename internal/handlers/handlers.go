package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dmawardi/bookings/internal/config"
	"github.com/dmawardi/bookings/internal/forms"
	"github.com/dmawardi/bookings/internal/models"
	"github.com/dmawardi/bookings/internal/render"
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

	render.AltRenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{
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

	render.AltRenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Renders Contact Page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["boo"] = "boodoo"
	render.AltRenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Renders Reservations Page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.ReservationForm
	formData := make(map[string]interface{})
	formData["reservation"] = emptyReservation

	stringMap := make(map[string]string)
	stringMap["boo"] = "boodoo"
	render.AltRenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		// Make form data empty
		Form: forms.NewForm(nil),
		Data: formData,
	})
}

// Handles posting of reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	// Parse form
	err := r.ParseForm()
	// Handle error
	if err != nil {
		log.Println("Error occurred while parsing form")
		return
	}

	// Extract form submitted values
	reservation := models.ReservationForm{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	// make a new form object
	form := forms.NewForm(r.PostForm)
	// Validation (Populates the Errors property)
	// Will check that form field is not an empty string
	form.Required("first_name", "last_name", "email")
	form.MinLength("last_name", 4)
	form.IsEmail("email")

	// if not valid
	if !form.Valid() {
		// Make map with reservation information
		data := make(map[string]interface{})
		// Pass input form values to data object
		data["reservation"] = reservation

		render.AltRenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			// Use form for retrieving errors
			Form: form,
			// Use data to return previously submitted values back to user
			Data: data,
		})
		return
	}

	// Put reservation data in Session
	m.App.Session.Put(r.Context(), "reservation", reservation)
	// Redirect user
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Renders Rooms: Majors page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["boo"] = "boodoo"
	render.AltRenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Renders Rooms: Generals page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["boo"] = "boodoo"
	render.AltRenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Renders Room Availability form page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["boo"] = "boodoo"
	render.AltRenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// POST Room Availability form
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	w.Write([]byte(fmt.Sprintf("start date is: %s and end date is: %s", start, end)))
}

// JSON response definition
type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// POST Availability route. Returns JSON with current status
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available",
	}

	// create Json using response (uses struct json details). No previx and indent with 5 spaces
	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		fmt.Println("Encountered error providing availability JSON")
	}

	// Edit content type
	w.Header().Set("Content-Type", "application/json")
	// Write data as response
	w.Write(out)
}

// Reservation Summary is shown after POSTing form
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	// reservation := m.App.Session.Get(r.Context(), "reservation")
	// ok is added to grab type casting error
	// Get reservation from session and cast as reservation form
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.ReservationForm)
	if !ok {
		log.Println("Couldn't find reservation in Sessions")
		// Place error in context sessions to display on page
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from Session")
		// redirect user (auto displays message upon redirect)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.AltRenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}
