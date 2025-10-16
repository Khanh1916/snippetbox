package main

import (
	"net/http"
	"testing"

	"github.com/Khanh1916/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	//create a new instance app contains logs
	app := newTestApplication(t)

	// create test server to test end-to-end
	testServer := newTestServer(t, app.routes()) // the request use all real app routes, middlewares, handlers
	defer testServer.Close()

	// catch a result of response from server when makeing a GET /ping
	code, _, body := testServer.get(t, "/ping")

	assert.Equal(t, code, http.StatusOK) //compare status code
	assert.Equal(t, body, "OK")          //compare body
}
