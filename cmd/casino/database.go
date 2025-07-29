package main

import (
	"context"
	"strings"
	"time"
)

// TODO: Transform this into a more elaborate query builder, surely
func buildFindTransactionsQueryWithLogicalAnds(userID string, typ TransactionType) (string, []any) {
	var (
		statement  strings.Builder
		positional = '1'
		arguments  = make([]any, 0, 2)
	)

	statement.WriteString(`SELECT
							   user_id,
							   transaction_type,
							   amount,
							   timestamp
						   FROM transactions`)

	if userID != "" {
		statement.WriteString("\nWHERE user_id = $")
		statement.WriteRune(positional)
		positional++

		arguments = append(arguments, userID)
	}

	if typ != "" {
		if positional == '1' {
			statement.WriteString("\nWHERE ")
		} else {
			statement.WriteString(" AND ")
		}
		statement.WriteString("transaction_type = $")
		statement.WriteRune(positional)
		positional++

		arguments = append(arguments, typ)
	}

	statement.WriteString(";")

	return statement.String(), arguments
}

func (app *app) findTransactionsFilterByType(ctx context.Context, typ TransactionType) ([]*Transaction, error) {
	return app.findTransactionsByUserID(ctx, "", typ)
}

func (app *app) findTransactionsByUserID(ctx context.Context, userID string, typ TransactionType) ([]*Transaction, error) {
	if typ == TransactionTypeAll {
		typ = ""
	}

	stmt, arguments := buildFindTransactionsQueryWithLogicalAnds(userID, typ)

	transactions := []*Transaction{}

	rows, err := app.db.Query(
		stmt,
		arguments...,
	)
	if err != nil {
		if strings.Contains(err.Error(), "22P02") {
			return transactions, nil
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			userID          int
			transactionType TransactionType
			amount          int
			timestamp       time.Time
		)

		t := &Transaction{
			UserID:          &userID,
			TransactionType: &transactionType,
			Amount:          &amount,
			Timestamp:       &timestamp,
		}

		if err := rows.Scan(t.UserID, t.TransactionType, t.Amount, t.Timestamp); err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
