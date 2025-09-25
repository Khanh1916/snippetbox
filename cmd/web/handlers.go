package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	// Gọi template
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.errorLog.Print(err.Error()) // log lỗi vào file log của ứng dụng thay vì log ra stdout
		//log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Render template
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.errorLog.Print(err.Error())
		//log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
	//w.Write([]byte("Toi dang test HTTP web"))
}

func (app *application) create(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("CREATING HTTP web"))
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
		w.Header().Set("Allow", "POST")
		// w.WriteHeader(405)
		// w.Write([]byte("Method Not Allowed"))
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Snippet created"))
}

func (app *application) jsonForTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	//w.Header()["CONTENT-TYPE"] = []string{"application/json"}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"name": "John", "age": 30}`))
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Displaying snippet with ID %d", id)
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
