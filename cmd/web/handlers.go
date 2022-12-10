package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"viktorkrams/snippetbox/pkg/models"
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
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

	id, err := app.snippets.Insert(title, text)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}

func (app *application) deleteSnippet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("deleting")
}
