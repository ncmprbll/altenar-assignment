package main_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	casino "github.com/ncmprbll/altenar-assignment/cmd/casino"
)

func TestNewConsumerProcessor(t *testing.T) {
	t.Run(fmt.Sprintf("TestCase_%02d", 0), func(t *testing.T) {
		processor, err := casino.NewTransactionProcessor(app.DB, 0, 1)
		if err != nil {
			t.Errorf("Failed to create transaction processor: %v", err)
		}

		consumer, err := casino.NewTransactionConsumer(app.Kafka, processor)
		if err != nil {
			t.Errorf("Failed to create transaction consumer: %v", err)
		}

		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		defer cancel()

		if err := consumer.Close(ctx); err != nil {
			t.Errorf("Failed to close transaction consumer: %v", err)
		}

		if err := processor.Close(ctx); err != nil {
			t.Errorf("Failed to close transaction processor: %v", err)
		}
	})
}

var transactions = []casino.Transaction{
	casino.NewTransaction(1, casino.TransactionTypeBet, 1, time.Now()),
	casino.NewTransaction(2, casino.TransactionTypeBet, 5, time.Now()),
	casino.NewTransaction(3, casino.TransactionTypeBet, 10, time.Now()),
	casino.NewTransaction(4, casino.TransactionTypeBet, 20, time.Now()),
	casino.NewTransaction(5, casino.TransactionTypeBet, 40, time.Now()),
	casino.NewTransaction(6, casino.TransactionTypeBet, 50, time.Now()),
	casino.NewTransaction(7, casino.TransactionTypeBet, 70, time.Now()),
	casino.NewTransaction(8, casino.TransactionTypeBet, 100, time.Now()),
	casino.NewTransaction(9, casino.TransactionTypeBet, 120, time.Now()),
	casino.NewTransaction(10, casino.TransactionTypeBet, 1000, time.Now()),
	casino.NewTransaction(11, casino.TransactionTypeBet, 1200, time.Now()),
}

func TestConsumerProcessorIntegration(t *testing.T) {
	t.Run(fmt.Sprintf("TestCase_%02d", 0), func(t *testing.T) {
		processor, err := casino.NewTransactionProcessor(app.DB, 0, 1)
		if err != nil {
			t.Errorf("Failed to create transaction processor: %v", err)
		}

		consumer, err := casino.NewTransactionConsumer(app.Kafka, processor)
		if err != nil {
			t.Errorf("Failed to create transaction consumer: %v", err)
		}

		for _, tx := range transactions {
			bytes, err := json.Marshal(tx)
			if err != nil {
				t.Errorf("Failed to marshal transaction: %v", err)
			}

			mock.ExpectExec("INSERT INTO transactions (user_id, transaction_type, amount, timestamp) VALUES ($1, $2, $3, $4);").
				WithArgs(tx.UserID, tx.TransactionType, tx.Amount, sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(1, 0))

			mockKafka.SendMessage(bytes)
		}

		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		defer cancel()

		if err := consumer.Close(ctx); err != nil {
			t.Errorf("Failed to close transaction consumer: %v", err)
		}

		if err := processor.Close(ctx); err != nil {
			t.Errorf("Failed to close transaction processor: %v", err)
		}
	})
}

var badTransactions = [][]byte{
	[]byte("{}"),
	[]byte("{\"user_id\":1}"),
	[]byte("not even a json object"),
	[]byte(`{"user_id":3,"transaction_type":"bet","amount":6,"timestamp":"2025-07-30T12:03:32+00:00"}`),
}

func TestConsumerProcessorIntegrationBad(t *testing.T) {
	t.Run(fmt.Sprintf("TestCase_%02d", 0), func(t *testing.T) {
		processor, err := casino.NewTransactionProcessor(app.DB, 0, 1)
		if err != nil {
			t.Errorf("Failed to create transaction processor: %v", err)
		}

		consumer, err := casino.NewTransactionConsumer(app.Kafka, processor)
		if err != nil {
			t.Errorf("Failed to create transaction consumer: %v", err)
		}

		for _, tx := range badTransactions {
			bytes, err := json.Marshal(tx)
			if err != nil {
				t.Errorf("Failed to marshal transaction: %v", err)
			}

			mockKafka.SendMessage(bytes)
		}

		ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
		defer cancel()

		if err := consumer.Close(ctx); err != nil {
			t.Errorf("Failed to close transaction consumer: %v", err)
		}

		if err := processor.Close(ctx); err != nil {
			t.Errorf("Failed to close transaction processor: %v", err)
		}
	})
}
