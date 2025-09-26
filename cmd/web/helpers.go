package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s/n%s", err.Error(), debug.Stack()) // debug.Stack() trả về stack trace
	app.errorLog.Output(2, trace)                              // ghi log với số dòng chính xác của lỗi (n=2)
	// http.Error(w, "Internal Server Error", 500)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
