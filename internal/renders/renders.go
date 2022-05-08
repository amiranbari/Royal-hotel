package renders

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/amiranbari/Royal-hotel/internal/config"
	"github.com/amiranbari/Royal-hotel/internal/models"
	"html/template"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

var functions = template.FuncMap{}

func NewRenderer(a *config.AppConfig) {
	app = a
}
func CreateTemplateCache() (config.TemplateCache, error) {
	myCache := config.TemplateCache{}

	pages, err := filepath.Glob("../../templates/*.page.tmpl")

	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("../../templates/*.layout.tmpl")

		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("../../templates/*.layout.tmpl")

			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	return td
}

func Template(rw http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc config.TemplateCache

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]

	if !ok {
		return errors.New("Could not get template from cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(rw)

	if err != nil {
		fmt.Println("Error writing template to browser", err)
		return nil
	}

	return nil
}
