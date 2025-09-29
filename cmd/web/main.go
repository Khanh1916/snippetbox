package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/Khanh1916/snippetbox/internal/models"
	_ "github.com/go-sql-driver/mysql"
)

// Cấu hình ứng dụng
type config struct {
	addr      string
	staticDir string
	dsn       string
}

// Phần phụ trợ của ứng dụng, inject các depedencies vào các handlers
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	cfg      config
	snippets *models.SnippetModel
}

func main() {
	// Khởi tạo cấu hình
	cfg := config{}

	// Đọc cờ dòng lệnh
	flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.StringVar(&cfg.dsn, "dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	// Tạo loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Mở kết nối đến MySQL
	db, err := openDB(cfg.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Tạo application
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		cfg:      cfg,
		snippets: &models.SnippetModel{DB: db},
	}

	// Tạo server HTTP để bắt được các error log từ server
	// (thay vì chỉ in ra stdout)
	srv := &http.Server{
		Addr:     app.cfg.addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	app.infoLog.Printf("Server started on %s", app.cfg.addr)
	err = srv.ListenAndServe()

	if err != nil {
		app.errorLog.Fatal(err)
	}
}

// Hàm openDB mở kết nối đến database, trả về sql.DB pool
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
