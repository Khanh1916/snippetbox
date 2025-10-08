package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	//"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Khanh1916/snippetbox/internal/models"
	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets
	// Use the new render helper.
	app.render(w, http.StatusOK, "home.html", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "create.html", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	// map includes errors in the validation
	fieldsErrors := make(map[string]string)

	if strings.TrimSpace(title) == "" {
		fieldsErrors["title"] = "This field cannot be blank."
	} else if utf8.RuneCountInString(title) > 100 {
		fieldsErrors["title"] = "This field cannot be more than 100 characters long."
	}

	if strings.TrimSpace(content) == "" {
		fieldsErrors["content"] = "This field cannot be blank."
	}

	if expires != 1 && expires != 7 && expires != 365 {
		fieldsErrors["expires"] = "This field must be equal 1, 7 or 365."
	}

	if len(fieldsErrors) > 0 {
		fmt.Fprint(w, fieldsErrors)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Lấy record từ database
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet
	// Use the new render helper.
	app.render(w, http.StatusOK, "view.html", data)
}

// // downloadFile phục vụ file tĩnh để tải về
// func (app *application) downloadFile(w http.ResponseWriter, r *http.Request) {
// 	filename := r.URL.Query().Get("file")

// 	// Chuẩn hóa đường dẫn
// 	cleanPath := filepath.Clean("./ui/static/" + filename)

// 	http.ServeFile(w, r, cleanPath)
// }

// noDirFileSystem là một wrapper quanh http.FileSystem để ngăn chặn việc liệt kê thư mục
type noDirFileSystem struct {
	fs http.FileSystem
}

// Open mở file và trả về lỗi nếu đó là một thư mục
func (n noDirFileSystem) Open(name string) (http.File, error) {
	f, err := n.fs.Open(name) // mở file
	if err != nil {
		return nil, err
	}
	s, err := f.Stat() // lấy thông tin file
	if err != nil {
		return nil, err
	}
	if s.IsDir() { // nếu là thư mục
		return nil, os.ErrNotExist
	}
	return f, nil
}
