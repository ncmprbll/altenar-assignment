package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	casino "github.com/ncmprbll/altenar-assignment/cmd/casino"
)

var getTransactionsTests = []struct {
	method            string
	url               string
	query             string
	expectedArgsCount int
	expectedCode      int
}{
	{
		method:            http.MethodGet,
		url:               "/transactions",
		query:             "?transaction_type=" + url.QueryEscape("bet;DROP TABLE users;"),
		expectedArgsCount: 1,
		expectedCode:      http.StatusOK,
	},
	{
		method:            http.MethodGet,
		url:               "/transactions",
		query:             "",
		expectedArgsCount: 0,
		expectedCode:      http.StatusOK,
	},
	{
		method:            http.MethodGet,
		url:               "/transactions",
		query:             "?transaction_type=" + string(casino.TransactionTypeBet),
		expectedArgsCount: 1,
		expectedCode:      http.StatusOK,
	},
	{
		method:            http.MethodGet,
		url:               "/transactions",
		query:             "?transaction_type=" + string(casino.TransactionTypeWin),
		expectedArgsCount: 1,
		expectedCode:      http.StatusOK,
	},
	{
		method:            http.MethodGet,
		url:               "/transactions",
		query:             "?transaction_type=" + string(casino.TransactionTypeAll),
		expectedArgsCount: 0,
		expectedCode:      http.StatusOK,
	},
	{
		method:            http.MethodGet,
		url:               "/transactions",
		query:             "?transaction_type=garbage",
		expectedArgsCount: 1,
		expectedCode:      http.StatusOK,
	},
}

func TestMethodGetTransactions(t *testing.T) {
	for i, tt := range getTransactionsTests {
		t.Run(fmt.Sprintf("TestCase_%02d", i), func(t *testing.T) {
			handler := app.Routes()

			// Ignore expectations
			switch tt.expectedArgsCount {
			case 0:
				mock.ExpectQuery("skip").WithoutArgs().WillReturnRows(sqlmock.NewRows([]string{}))
			case 1:
				mock.ExpectQuery("skip").WithArgs(sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{}))
			}

			target := tt.url + tt.query
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.method, target, nil)

			handler.ServeHTTP(w, r)
			statusCode := w.Result().StatusCode
			if statusCode != tt.expectedCode {
				t.Errorf("%s %s status code mismatch: got %d, want %d", tt.method, target, statusCode, tt.expectedCode)
			}
		})
	}
}

var getUserTransactionsTests = []struct {
	method            string
	url               string
	query             string
	expectedArgsCount int
	expectedCode      int
}{
	{
		method:            http.MethodGet,
		url:               "/users/1/transactions",
		query:             "?transaction_type=" + url.QueryEscape("bet;DROP TABLE users;"),
		expectedArgsCount: 2,
		expectedCode:      http.StatusOK,
	},
	{
		method:            http.MethodGet,
		url:               "/users/" + url.QueryEscape("1;DROP TABLE users;") + "/transactions",
		expectedArgsCount: 1,
		expectedCode:      http.StatusOK,
	},
	{
		method:            http.MethodGet,
		url:               "/users/1/transactions",
		expectedArgsCount: 1,
		expectedCode:      http.StatusOK,
	},
	{
		method:            http.MethodGet,
		url:               "/users/1/transactions",
		query:             "?transaction_type=" + string(casino.TransactionTypeBet),
		expectedArgsCount: 2,
		expectedCode:      http.StatusOK,
	},
	{
		method:            http.MethodGet,
		url:               "/users/1/transactions",
		query:             "?transaction_type=" + string(casino.TransactionTypeWin),
		expectedArgsCount: 2,
		expectedCode:      http.StatusOK,
	},
	{
		method:            http.MethodGet,
		url:               "/users/1/transactions",
		query:             "?transaction_type=" + string(casino.TransactionTypeAll),
		expectedArgsCount: 1,
		expectedCode:      http.StatusOK,
	},
}

func TestMethodGetUserTransactions(t *testing.T) {
	for i, tt := range getUserTransactionsTests {
		t.Run(fmt.Sprintf("TestCase_%02d", i), func(t *testing.T) {
			handler := app.Routes()

			// Ignore expectations
			switch tt.expectedArgsCount {
			case 1:
				mock.ExpectQuery("skip").WithArgs(sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{}))
			case 2:
				mock.ExpectQuery("skip").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{}))
			}

			target := tt.url + tt.query
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.method, target, nil)

			handler.ServeHTTP(w, r)
			statusCode := w.Result().StatusCode
			if statusCode != tt.expectedCode {
				t.Errorf("%s %s status code mismatch: got %d, want %d", tt.method, target, statusCode, tt.expectedCode)
			}
		})
	}
}

func TestMethodGetTransactionsInternalServerError(t *testing.T) {
	t.Run(fmt.Sprintf("TestCase_%02d", 0), func(t *testing.T) {
		handler := app.Routes()

		target := "/transactions?transaction_type=type"
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, target, nil)

		handler.ServeHTTP(w, r)
		statusCode := w.Result().StatusCode
		if statusCode != http.StatusInternalServerError {
			t.Errorf("%s %s status code mismatch: got %d, want %d", http.MethodGet, target, statusCode, http.StatusInternalServerError)
		}
	})
}

func TestMethodGetUserTransactionsInternalServerError(t *testing.T) {
	t.Run(fmt.Sprintf("TestCase_%02d", 0), func(t *testing.T) {
		handler := app.Routes()

		target := "/users/1/transactions?transaction_type=type"
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, target, nil)

		handler.ServeHTTP(w, r)
		statusCode := w.Result().StatusCode
		if statusCode != http.StatusInternalServerError {
			t.Errorf("%s %s status code mismatch: got %d, want %d", http.MethodGet, target, statusCode, http.StatusInternalServerError)
		}
	})
}
