package main

import "time"

type TransactionType string

const (
	TransactionBet TransactionType = "bet"
	TransactionWin TransactionType = "win"
)

type Transaction struct {
	UserID          *int             `json:"user_id"`
	TransactionType *TransactionType `json:"transaction_type"`
	Amount          *int             `json:"amount"`
	Timestamp       *time.Time       `json:"timestamp"`
}

func (t Transaction) HasAllFields() bool {
	return t.UserID != nil && t.TransactionType != nil && t.Amount != nil && t.Timestamp != nil
}

func (t Transaction) Valid() bool {
	return t.HasAllFields() && *t.UserID > 0
}
