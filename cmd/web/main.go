package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// Cấu hình ứng dụng
type config struct {
	addr      string
	staticDir string
}

// Phần phụ trợ của ứng dụng, inject các depedencies vào các handlers
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Khởi tạo cấu hình ứng dụng
	var cfg config

	// Đọc cờ dòng lệnh
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Tạo application
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux := http.NewServeMux()
	// Đăng các route cụ thể trước, rồi đến catch-all
	//mux.HandleFunc("/view/", view)    // subtree: /view/...
	mux.HandleFunc("/view", app.snippetView) // fixed path: /view
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	mux.HandleFunc("example.com/local", app.handlerHost) // host-specific
	mux.HandleFunc("/local", app.handlerGeneral)         // non-host-specific

	mux.HandleFunc("/create", app.create) // fixed path: /create
	mux.HandleFunc("/json", app.jsonForTest)

	// File server cho static assets
	fs := noDirFileSystem{http.Dir(cfg.staticDir)}
	// StripPrefix để bỏ "/static" khỏi URL trước khi gửi đến file server
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(fs)))

	// Catch-all cuối cùng
	mux.HandleFunc("/", app.home)

	// Tạo server HTTP để bắt được các error log từ server
	// (thay vì chỉ in ra stdout)
	srv := &http.Server{
		Addr:     cfg.addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Server started on %s", cfg.addr)
	err := srv.ListenAndServe()

	errorLog.Fatal(err)
}
