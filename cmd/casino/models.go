package main

import (
	"time"
)

type TransactionType string

const (
	TransactionTypeBet TransactionType = "bet"
	TransactionTypein  TransactionType = "win"
	TransactionTypeAll TransactionType = "all"
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
	return t.HasAllFields() && *t.UserID > 0 && *t.Amount > 0
}
