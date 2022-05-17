package main

import (
	"fmt"
	"github.com/amiranbari/Royal-hotel/internal/config"
	"github.com/go-chi/chi/v5"
	"testing"
)

func TestRoute(t *testing.T) {
	var app config.AppConfig

	h := route(&app)

	switch v := h.(type) {
	case *chi.Mux:
		//do nothing
	default:
		t.Error(fmt.Sprintf("type is not http.handler, but is %T", v))
	}
}
