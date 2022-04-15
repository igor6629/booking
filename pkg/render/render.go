package render

import (
	"bytes"
	"fmt"
	"github.com/igor6629/booking/pkg/config"
	"github.com/igor6629/booking/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}
var app *config.AppConfig

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func NewTemplates(conf *config.AppConfig) {
	app = conf
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]

	if !ok {
		log.Fatal("Could not get template from template")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td)
	_ = t.Execute(buf, td)
	_, err := buf.WriteTo(w)

	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./templates/*.page.tmpl")

	if err != nil {
		return cache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return cache, err
		}

		matches, err := filepath.Glob("./templates/*layout.tmpl")

		if err != nil {
			return cache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*layout.tmpl")

			if err != nil {
				return cache, err
			}
		}

		cache[name] = ts
	}

	return cache, nil
}
