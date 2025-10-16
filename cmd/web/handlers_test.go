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

func TestSnippetView(t *testing.T) {
	app := newTestApplication(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/snippet/view/1",
			wantCode: http.StatusOK,
			wantBody: "An old silent pond...",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/snippet/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/snippet/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/snippet/view/1.23",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/snippet/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/snippet/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.urlPath)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}
