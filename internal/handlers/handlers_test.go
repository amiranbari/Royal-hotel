package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var theTests = []struct {
	name               string
	url                string
	expectedStatusCode int
}{
	{"index", "/", http.StatusOK},
	{"search", "/search", http.StatusOK},
	{"login", "/login", http.StatusOK},
	{"logout", "/logout", http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()

	ts := httptest.NewTLSServer(routes)

	defer ts.Close()

	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}

		//} else {
		//
		//	values := url.Values{}
		//	for _, x := range e.params {
		//		values.Add(x.key, x.value)
		//	}
		//
		//	resp, err := ts.Client().PostForm(ts.URL+e.url, values)
		//
		//	if err != nil {
		//		t.Log(err)
		//		t.Fatal(err)
		//	}
		//
		//	if resp.StatusCode != e.expectedStatusCode {
		//		t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		//	}
	}
}
