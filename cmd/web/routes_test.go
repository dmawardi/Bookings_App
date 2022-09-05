package main

import (
	"fmt"
	"testing"

	"github.com/dmawardi/bookings/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	// Generate app for use in routes
	var app config.AppConfig

	// Use routes to build router in routes
	mux := routes(&app)

	switch v := mux.(type) {
	// If the mux router from CHI
	case *chi.Mux:
		// Consider pass, do nothing

	default:
		t.Error(fmt.Sprintf("type is %T, expected Http Handler", v))
	}

}
