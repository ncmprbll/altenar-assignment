package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *app) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.StripSlashes)

	r.Get("/transactions", app.getTransactions)
	r.Get("/users/{userID}/transactions", app.getUserTransactions)

	return r
}
