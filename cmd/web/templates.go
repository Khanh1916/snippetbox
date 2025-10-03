package main

import (
	"html/template"
	"path/filepath"

	"github.com/Khanh1916/snippetbox/internal/models"
)

// struct lưu trữ được nhiều mảnh dynamic data thay vì chỉ một mảnh
type templateData struct {
	CurrentYear int // Add a CurrentYear field to the templateData struct.
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

// cache template avoid duplicate parsing files many times
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
