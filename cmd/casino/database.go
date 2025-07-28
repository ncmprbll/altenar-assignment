package main

import "context"

func (app *app) findTransactionsFilterByType(ctx context.Context, typ TransactionType) ([]*Transaction, error) {
	return []*Transaction{}, nil
}

func (app *app) findTransactionsByUserID(ctx context.Context, userID string, typ TransactionType) ([]*Transaction, error) {
	return []*Transaction{}, nil
}
