package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() *chi.Mux {
	mux := chi.NewRouter()

	mux.Route("/", func(r chi.Router) {
		r.Get("/", app.home)

		r.Route("/snippet", func(r chi.Router) {
			r.Get("/{id}", app.showSnippet)

			r.Route("/create", func(r chi.Router) {
				r.Get("/", app.createSnippet)
			})

			r.Route("/post", func(r chi.Router) {
				r.Post("/", app.addSnippet)
			})

			r.Route("/delete", func(r chi.Router) {
				r.Get("/{id}", app.deleteSnippet)
			})
		})
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return mux
}
