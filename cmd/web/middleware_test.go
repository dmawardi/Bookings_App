package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myHandler MyHandler
	h := NoSurfCSRF(&myHandler)

	// Find type of result
	switch v := h.(type) {
	// If http handler, consider as pass
	case http.Handler:
		// do nothing
	default:
		// %T is for type
		t.Error(fmt.Sprintf("type is %T, expected Http Handler", v))
	}
}

func TestSessionLoad(t *testing.T) {
	var myHandler MyHandler
	h := SessionLoad(&myHandler)

	switch v := h.(type) {
	case http.Handler:
	// If http handler, consider as pass
	default:
		// %T is for type
		t.Error(fmt.Sprintf("type is %T, expected Http Handler", v))
	}
}
