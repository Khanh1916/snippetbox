package main

import "net/http"

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()
	// Đăng các route cụ thể trước, rồi đến catch-all
	//mux.HandleFunc("/view/", view)    // subtree: /view/...
	mux.HandleFunc("/snippet/view", app.snippetView) // fixed path: /view
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	mux.HandleFunc("example.com/local", app.handlerHost) // host-specific
	mux.HandleFunc("/local", app.handlerGeneral)         // non-host-specific

	// File server cho static assets
	fs := noDirFileSystem{http.Dir(app.cfg.staticDir)}
	// StripPrefix để bỏ "/static" khỏi URL trước khi gửi đến file server
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(fs)))

	// Catch-all cuối cùng
	mux.HandleFunc("/", app.home)
	return secureHeaders(mux)
}
