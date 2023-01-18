package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/NPeykov/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
    app := newTestApplication(t)
    ts := newTestServer(t, app.routes())
    defer ts.Close()

    code, _, body := ts.get(t, "/ping")

    assert.Equal(t, code, http.StatusOK)
    assert.Equal(t, body, "OK")
}

func TestSnippetView(t *testing.T) {
    app := newTestApplication(t)
    ts := newTestServer(t, app.routes())
    defer ts.Close()

    tests := []struct {
        name string
        urlPath string
        wantCode int
        wantBody string
    }{
        {
            name: "Valid ID",
            urlPath: "/snippet/view/1",
            wantCode: http.StatusOK,
            wantBody: "The art of saying art",
        },
        {
            name: "Non-existent ID",
            urlPath: "/snippet/view/2",
            wantCode: http.StatusNotFound,
            wantBody: "",
        },
        {
            name: "Negative ID",
            urlPath: "/snippet/view/-1",
            wantCode: http.StatusNotFound,
            wantBody: "",
        },
        {
            name: "Decimal ID",
            urlPath: "/snippet/view/1.23",
            wantCode: http.StatusNotFound,
            wantBody: "",
        },
        {
            name: "String ID",
            urlPath: "/snippet/view/foo",
            wantCode: http.StatusNotFound,
            wantBody: "",
        },
        {
            name: "Empty ID",
            urlPath: "/snippet/view/",
            wantCode: http.StatusNotFound,
            wantBody: "",
        },
    }

    for _, test := range tests { 
        t.Run(test.name, func(t *testing.T) {
            code, _, body := ts.get(t, test.urlPath)
            assert.Equal(t, code, test.wantCode)
            assert.StringContains(t, body, test.wantBody)
        })
    }
}

func TestUserSignup(t *testing.T) {
    app := newTestApplication(t)
    ts := newTestServer(t, app.routes())
    defer ts.Close()

    _, _, body := ts.get(t, "/user/signup")
    validCSRFToken := extractCSRFToken(t, body)

    const (
        validName = "pkv"
        validPassword = "12345678"
        validEmail = "WP@gg.com"
        formTag = `<form action="/user/signup" method="POST" novalidate>`
    )

    tests := []struct {
        name string
        userName string
        userEmail string
        userPassword string
        csrfToken string
        wantCode int
        wantFormTag string
    }{
        {
            name: "Valid submission",
            userName: validName,
            userEmail: validEmail,
            userPassword: validPassword,
            csrfToken: validCSRFToken,
            wantCode: http.StatusSeeOther,
            wantFormTag: "",
        },
        {
            name: "Invalid CSRF Token",
            userName: validName,
            userEmail: validEmail,
            userPassword: validPassword,
            csrfToken: "wrongToken",
            wantCode: http.StatusBadRequest,
            wantFormTag: "",
        },
        {
            name: "Empty name",
            userName: "",
            userEmail: validEmail,
            userPassword: validPassword,
            csrfToken: validCSRFToken,
            wantCode: http.StatusUnprocessableEntity,
            wantFormTag: formTag,
        },
        {
            name: "Empty email",
            userName: validName,
            userEmail: "",
            userPassword: validPassword,
            csrfToken: validCSRFToken,
            wantCode: http.StatusUnprocessableEntity,
            wantFormTag: formTag,
        },
        {
            name: "Empty password",
            userName: validName,
            userEmail: validEmail, 
            userPassword: "",
            csrfToken: validCSRFToken,
            wantCode: http.StatusUnprocessableEntity,
            wantFormTag: formTag,
        },
        {
            name: "Invalid Email",
            userName: validName,
            userEmail: "myEmail@", 
            userPassword: validPassword,
            csrfToken: validCSRFToken,
            wantCode: http.StatusUnprocessableEntity,
            wantFormTag: formTag,
        },
        {
            name: "Short password",
            userName: validName,
            userEmail: validEmail, 
            userPassword: "123",
            csrfToken: validCSRFToken,
            wantCode: http.StatusUnprocessableEntity,
            wantFormTag: formTag,
        },
        {
            name: "Duplicate email",
            userName: validName,
            userEmail: "gg@gg.com", 
            userPassword: validPassword,
            csrfToken: validCSRFToken,
            wantCode: http.StatusUnprocessableEntity,
            wantFormTag: formTag,
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            form := url.Values{}
            form.Add("name", test.userName)
            form.Add("email", test.userEmail)
            form.Add("password", test.userPassword)
            form.Add("csrf_token", test.csrfToken)

            code, _, body := ts.postForm(t, "/user/signup", form)

            assert.Equal(t, code, test.wantCode)
            assert.StringContains(t, body, test.wantFormTag)
        })
    }
}


