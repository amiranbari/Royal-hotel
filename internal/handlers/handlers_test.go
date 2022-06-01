package handlers

import (
	"context"
	"log"
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
	}
}

func Test_SearchForRoom(t *testing.T) {
	//test without body request
	req, _ := http.NewRequest("GET", "/search", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.SearchForRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("SearchForRoom Handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	//test with wrong start date
	req, _ = http.NewRequest("GET", "/search?start_date=wronStartDate&end_date=2050-01-01", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("SearchForRoom Handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusInternalServerError)
	}

	//test with wrong end date
	req, _ = http.NewRequest("GET", "/search?start_date=2050-01-01&end_date=wrongEndDate", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("SearchForRoom Handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusInternalServerError)
	}

	//test with error in database searching
	req, _ = http.NewRequest("GET", "/search?start_date=2050-01-01&end_date=2050-01-03", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("SearchForRoom Handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test with non room finding
	req, _ = http.NewRequest("GET", "/search?start_date=2050-01-01&end_date=2050-01-04", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("SearchForRoom Handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test with everything ok
	req, _ = http.NewRequest("GET", "/search?start_date=2050-01-01&end_date=2050-01-05", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("SearchForRoom Handler return wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

}

func getCtx(r *http.Request) context.Context {
	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
