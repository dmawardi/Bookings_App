package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/dmawardi/bookings/internal/config"
)

var app *config.AppConfig

// Function called in main.go to connect app state to current file
func SetState(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of: ", status)
	// Make HTTP error
	http.Error(w, http.StatusText(status), status)

}

func ServerError(w http.ResponseWriter, err error) {
	// Build debug trace string using error and debug stack
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
