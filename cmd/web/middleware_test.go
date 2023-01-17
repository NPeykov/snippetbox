package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NPeykov/snippetbox/internal/assert"
)

func TestSecureHeaders(t *testing.T) {
    recorder := httptest.NewRecorder()

    req, err := http.NewRequest(http.MethodGet, "/", nil)
    if err != nil {
        t.Fatal(err)
    }

    next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("I'm after the middleware"))
    })

    secureHeaders(next).ServeHTTP(recorder, req)
    rs := recorder.Result()

    headersToTest := []struct {
        headerKey string
        expectedHeaderValue string
    }{
        {
            headerKey: "Content-Security-Policy",
            expectedHeaderValue: "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com",
        },
        {
            headerKey: "Referrer-Policy",
            expectedHeaderValue: "origin-when-cross-origin",
        },
        {
            headerKey: "X-Content-Type-Options",
            expectedHeaderValue: "nosniff",
        },
        {
            headerKey: "X-Frame-Options",
            expectedHeaderValue: "deny",
        },
        {
            headerKey: "X-XSS-Protection",
            expectedHeaderValue: "0",
        },
    }
    
    for _, test := range headersToTest {
        assert.Equal(t, rs.Header.Get(test.headerKey), test.expectedHeaderValue)
    }

    assert.Equal(t, rs.StatusCode, http.StatusOK)

    defer rs.Body.Close()
    body, err := io.ReadAll(rs.Body)
    if err != nil {
        t.Fatal(err)
    }

    bytes.TrimSpace(body)
    assert.Equal(t, string(body), "I'm after the middleware")
}
