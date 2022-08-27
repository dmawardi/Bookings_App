package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func WriteToConsole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hit the page: ", r.URL)
		// Proceed to the next middleware with w and r
		next.ServeHTTP(w, r)
		fmt.Println("Leaving the server: ")

	})
}

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
