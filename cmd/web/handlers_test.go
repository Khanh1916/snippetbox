package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Khanh1916/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	//create a new instance app contains logs
	app := &application{
		errorLog: log.New(io.Discard, "", 0),
		infoLog:  log.New(io.Discard, "", 0),
	}

	// create test server to test end-to-end
	testServer := httptest.NewTLSServer(app.routes()) // the request use all real app routes, middlewares, handlers
	defer testServer.Close()

	// catch a result of response from server when makeing a GET /ping
	rs, err := testServer.Client().Get(testServer.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, rs.StatusCode, http.StatusOK) //compare status code

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK") //compare body
}
