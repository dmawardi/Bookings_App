package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/dmawardi/bookings/internal/config"
	"github.com/dmawardi/bookings/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	gob.Register(models.ReservationForm{})

	// change this to true when in production
	testApp.InProduction = false

	// Create info log that outputs to std output, has prefix INFO, then date/time
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	// Set app info log as created logger
	testApp.InfoLog = infoLog

	// Create error log that outputs to std output, has prefix INFO, then date/time
	// short file provides information on error
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	// Set app error log as created logger
	testApp.ErrorLog = errorLog

	// set up the session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}

func (tw *myWriter) WriteHeader(i int) {}

func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
