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
	cfg      config
}

func main() {
	// Khởi tạo cấu hình
	cfg := config{}

	// Tạo loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Tạo application
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		cfg:      cfg,
	}
	// Đọc cờ dòng lệnh
	flag.StringVar(&app.cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&app.cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()

	// Tạo server HTTP để bắt được các error log từ server
	// (thay vì chỉ in ra stdout)
	srv := &http.Server{
		Addr:     app.cfg.addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	app.infoLog.Printf("Server started on %s", app.cfg.addr)
	err := srv.ListenAndServe()

	app.errorLog.Fatal(err)
}
