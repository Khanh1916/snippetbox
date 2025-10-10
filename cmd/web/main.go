package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Khanh1916/snippetbox/internal/models"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
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
	errorLog       *log.Logger
	infoLog        *log.Logger
	cfg            config
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template //add templateCache to application struct
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager //add session manager for flash messages
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

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true //cookie only sent through https with TLS

	// Tạo application
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		cfg:            cfg,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.CurveP256, tls.X25519},
	}

	// Tạo server HTTP để bắt được các error log từ server
	// (thay vì chỉ in ra stdout) cấu hình ở phần tạo logger
	srv := &http.Server{
		Addr:      app.cfg.addr,
		ErrorLog:  errorLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
		// add Idle, Read, Write timeouts to the server
		IdleTimeout:  time.Minute, //Idle should be set after setting ReadTimeout
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second, //write > read
	}

	app.infoLog.Printf("Server started on %s", app.cfg.addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")

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
