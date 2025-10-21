package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/Khanh1916/snippetbox/internal/models"
	"github.com/Khanh1916/snippetbox/ui"
)

// struct lưu trữ được nhiều mảnh dynamic data thay vì chỉ một mảnh
type templateData struct {
	CurrentYear     int // Add a CurrentYear field to the templateData struct
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any    // Add Form to refill user data after error validating and displaying error in html
	Flash           string // Add a Flash field to display flash message
	IsAuthenticated bool
	CSRFToken       string // Add CSRFToken field
	User            *models.User
}

// cache template avoid duplicate parsing files many times
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func humanDate(t time.Time) string {
	// Return the empty string if time has the zero value.
	if t.IsZero() {
		return ""
	}
	// Convert the time to UTC before formatting it.
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}
