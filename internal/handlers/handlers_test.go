package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	// GET
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"gov", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"ms", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"search", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"reserve", "/make-reservation", "GET", []postData{}, http.StatusOK},
	// POST
	{"POSTSearch", "/search-availability", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-03"},
	}, http.StatusOK},
	{"POSTSearchJSON", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2020-01-01"},
		{key: "end", value: "2020-01-03"},
	}, http.StatusOK},
	{"POSTReserve", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "Jim"},
		{key: "last_name", value: "Rodgers"},
		{key: "phone", value: "510-666-343"},
		{key: "email", value: "abdul@jabar.com"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)

	// Close test server when Test handlers function complerted
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			// GET request to test server URL + URL to test
			resp, err := ts.Client().Get(ts.URL + e.url)

			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			// Ensure response is ok
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf(`for %s, expected %d but got %d`, e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			// url values is a built in type that holds info as a POST request
			values := url.Values{}

			// Build data to post
			for _, POSTparameter := range e.params {
				values.Add(POSTparameter.key, POSTparameter.value)
			}

			// POST request
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			// Ensure response is ok
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf(`for %s, expected %d but got %d`, e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
