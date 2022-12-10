package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"viktorkrams/snippetbox/pkg/models"

	"github.com/go-chi/chi/v5"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.html", &templateData{Snippets: snippets})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
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

	app.render(w, r, "show.page.html", &templateData{Snippet: snippet})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	snippet := models.Snippet{}

	app.render(w, r, "create.page.html", &templateData{Snippet: &snippet})
}

func (app *application) addSnippet(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("Title")
	text := r.FormValue("Text")

	if title == "" {
		if len(text) > 10 {
			title = text[:10] + "..."
		} else {
			title = text
		}
	}

	id, err := app.snippets.Insert(title, text)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) deleteSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	err = app.snippets.Delete(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/"), http.StatusSeeOther)
}
