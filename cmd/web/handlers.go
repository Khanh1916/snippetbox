package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Khanh1916/snippetbox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Use the new render helper.
	app.render(w, http.StatusOK, "home.html", &templateData{Snippets: snippets})
}

// Handler cho host-specific
func (app *application) handlerHost(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("Host-specific handler for host: %s", r.Host)))
}

// Handler cho non-host-specific
func (app *application) handlerGeneral(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf("General handler (non-host-specific)")))
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "0 snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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
	// Use the new render helper.
	app.render(w, http.StatusOK, "view.html", &templateData{Snippet: snippet})
}

// downloadFile phục vụ file tĩnh để tải về
func (app *application) downloadFile(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("file")

	// Chuẩn hóa đường dẫn
	cleanPath := filepath.Clean("./ui/static/" + filename)

	http.ServeFile(w, r, cleanPath)
}

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
