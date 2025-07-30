package main

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/go-chi/chi/v5"
)

func internalServerError(w http.ResponseWriter, err error) {
	log.Printf("%s\n%s", err.Error(), debug.Stack())
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *App) GetTransactions(w http.ResponseWriter, r *http.Request) {
	typ := TransactionType(r.URL.Query().Get("transaction_type"))

	transactions, err := app.FindTransactionsFilterByType(r.Context(), typ)
	if err != nil {
		internalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		internalServerError(w, err)
	}
}

func (app *App) GetUserTransactions(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	typ := TransactionType(r.URL.Query().Get("transaction_type"))

	transactions, err := app.FindTransactionsByUserID(r.Context(), userID, typ)
	if err != nil {
		internalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		internalServerError(w, err)
	}
}
