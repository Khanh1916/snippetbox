package main

import (
	"github.com/Khanh1916/snippetbox/internal/models"
)

// struct lưu trữ được nhiều mảnh dynamic data thay vì chỉ một mảnh
type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
