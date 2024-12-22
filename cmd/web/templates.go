package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/ashab-k/snippetbox/pkg/forms"
	"github.com/ashab-k/snippetbox/pkg/models"
)

type templateData struct {
	CurrentYear int
	Form *forms.Form
	Flash string
	Snippet *models.Snippet
	Snippets []*models.Snippet
}

func humanDate(t time.Time) string {
	// Return the empty string if time has the zero value.
	if t.IsZero() {
		return ""
	}

	// Convert the time to UTC before formatting it.
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

// initialise a template.FuncMap object to store a lookup reference
// to all our self created template functions
var functions = template.FuncMap{
	"humanDate": humanDate,
}


func newTemplateCache(dir string) (map[string]*template.Template, error) {
    cache := map[string]*template.Template{}

    pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
    if err != nil {
        return nil, err
    }

    for _, page := range pages {
        name := filepath.Base(page)

        ts, err := template.New(name).Funcs(functions).ParseFiles(page)
        if err != nil {
            return nil, err
        }

        ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
        if err != nil {
            return nil, err
        }

        ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
        if err != nil {
            return nil, err
        }

        cache[name] = ts
    }

    return cache, nil
}