package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Khanh1916/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	// create Recorder
	rr := httptest.NewRecorder()

	// Initialize a new dummy http.Request.
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	//call ping to pass Recorder and Request
	ping(rr, r)

	rs := rr.Result() // got a result

	assert.Equal(t, rs.StatusCode, http.StatusOK) //compare status code

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK") //compare body
}
