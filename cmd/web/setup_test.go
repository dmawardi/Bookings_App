package main

import (
	"net/http"
	"os"
	"testing"
)

// TestMain sets up the environment before tests
func TestMain(m *testing.M) {
	// Complete before tests

	// Exit main run
	os.Exit(m.Run())
}

type MyHandler struct {
}

func (mh *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
