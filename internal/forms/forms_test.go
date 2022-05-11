package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func Test_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	form := New(r.PostForm)
	if !form.Valid() {
		t.Error("Something error in form valid.")
	}
}

func Test_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	form := New(r.PostForm)
	form.Required("whatever")
	if form.Valid() {
		t.Error("Form should not be valid!")
	}

	postedData := url.Values{}
	postedData.Add("whatever", "whatever")
	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("whatever")

	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func Test_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	form := New(r.PostForm)
	if form.Has("whatever") {
		t.Error("Form should not be valid!")
	}

	postedData := url.Values{}
	postedData.Add("whatever", "whatever")
	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)

	if !form.Has("whatever") {
		t.Error("shows does not have fields when it does")
	}
}

func Test_IsEmail(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	form := New(r.PostForm)
	if form.IsEmail("whatever") {
		t.Error("Form does not have email value")
	}

	postedData := url.Values{}
	postedData.Add("whatever", "whatever@gmail.com")
	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)

	if !form.IsEmail("whatever") {
		t.Error("shows it is not an email when it is")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	postedData.Add("whatever", "a")
	r := httptest.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData

	form := New(r.PostForm)
	if form.MinLength("whatever", 2) {
		t.Error("shows we have min-length when we dont")
	}

	postedData = url.Values{}
	postedData.Add("whatever", "ab")
	r, _ = http.NewRequest("POST", "/whatever", nil)

	r.PostForm = postedData
	form = New(r.PostForm)

	if !form.MinLength("whatever", 2) {
		t.Error("shows we dont have min-length when we have")
	}
}
