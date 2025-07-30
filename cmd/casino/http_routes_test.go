package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestMethodRoutes(t *testing.T) {
	t.Run(fmt.Sprintf("TestCase_%02d", 0), func(t *testing.T) {
		handler := app.Routes()

		mock.ExpectQuery("skip").WithoutArgs().WillReturnRows(sqlmock.NewRows([]string{}))

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/transactions", nil)

		handler.ServeHTTP(w, r)
		if w.Result().StatusCode == http.StatusNotFound {
			t.Errorf("route /transactions doesn't exist")
		}

		mock.ExpectQuery("skip").WithArgs(sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{}))

		w = httptest.NewRecorder()
		r = httptest.NewRequest(http.MethodGet, "/users/1/transactions", nil)

		handler.ServeHTTP(w, r)
		if w.Result().StatusCode == http.StatusNotFound {
			t.Errorf("route /users/{userID}/transactions doesn't exist")
		}
	})
}
