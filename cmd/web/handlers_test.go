package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NPeykov/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
    recorder := httptest.NewRecorder()

    req, err := http.NewRequest(http.MethodGet, "/", nil)
    if err != nil {
        t.Fatal(err)
    }

    ping(recorder, req)

    rs := recorder.Result()

    assert.Equal(t, rs.StatusCode, http.StatusOK)

    defer rs.Body.Close()
    body, err := io.ReadAll(rs.Body)
    if err != nil {
        t.Fatal(err)
    }

    bytes.TrimSpace(body)
    assert.Equal(t, string(body), "OK")
}
