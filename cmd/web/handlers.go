package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/NPeykov/snippetbox/internal/models"
	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
    snippets, err := app.snippets.Latest()

    if err != nil {
        app.serverError(w, err)
        return
    }

    tmplData := app.newTemplateData()
    tmplData.Snippets = snippets

    app.render(w, http.StatusOK, "home.tmpl.html", tmplData)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
    params := httprouter.ParamsFromContext(r.Context())
    id, err := strconv.Atoi(params.ByName("id"))
    if err != nil || id < 0 {
        app.notFound(w)
        return
    }
    snippet, err := app.snippets.Get(id)
    if err != nil {
        if errors.Is(err, models.ErrNoRecord) {
            app.notFound(w)
        } else {
            app.serverError(w, err)
        }
        return
    }
    
    tmplData := app.newTemplateData()
    tmplData.Snippet = snippet

    app.render(w, http.StatusOK, "view.tmpl.html", tmplData)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Display the form for creating a new snippet..."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
    title := "Everybody is changing"
    content := "People and time aswell are changing"
    expires := 7

    id, err := app.snippets.Insert(title, content, expires)

    if err != nil {
        app.serverError(w, err)
    }

    http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

