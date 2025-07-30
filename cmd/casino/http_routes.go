package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *App) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.StripSlashes)

	r.Get("/transactions", app.GetTransactions)
	r.Get("/users/{userID}/transactions", app.GetUserTransactions)

	return r
}
