package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.StripSlashes)

	r.Get("/users/{userID}/transactions", func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "userID")

		w.Write([]byte("HELLO, " + userID))
	})

	r.Get("/transactions", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HELLO, /transactions"))
	})

	return r
}
