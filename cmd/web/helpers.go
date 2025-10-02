package main

import (
	"bytes"
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

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	//Khoi tao buffer de test thu runtime error
	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err) //báo biến w và trả về error
	}

	w.WriteHeader(status) //Neu thanh cong se tra ve 200

	buf.WriteTo(w) //Viet ket qua cua buffer vao w
}
