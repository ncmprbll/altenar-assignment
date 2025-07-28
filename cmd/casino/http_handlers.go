package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func internalServerError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *app) getUserTransactions(w http.ResponseWriter, r *http.Request) {
	typ := TransactionType(r.URL.Query().Get("transaction_type"))

	transactions, err := app.findTransactionsFilterByType(r.Context(), typ)
	if err != nil {
		internalServerError(w)
		return
	}

	if json.NewEncoder(w).Encode(transactions); err != nil {
		internalServerError(w)
	}
}

func (app *app) getTransactions(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	typ := TransactionType(r.URL.Query().Get("transaction_type"))

	transactions, err := app.findTransactionsByUserID(r.Context(), userID, typ)
	if err != nil {
		internalServerError(w)
		return
	}

	if json.NewEncoder(w).Encode(transactions); err != nil {
		internalServerError(w)
	}
}
