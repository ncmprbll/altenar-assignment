package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type transactionProcessor struct {
	pgPool         *pgxpool.Pool
	transactionsCh chan *Transaction

	shutdownContext context.Context
	shutdownCancel  context.CancelFunc
	wg              sync.WaitGroup
}

// Initialize a new transaction processor
func NewTransactionProcessor(pgPool *pgxpool.Pool, jobsSize, workersCount int) *transactionProcessor {
	if jobsSize < 0 {
		panic("jobsSize cannot be negative")
	}

	if workersCount <= 0 {
		panic("workersCount cannot be negative or zero")
	}

	ctx, cancel := context.WithCancel(context.Background())

	processor := &transactionProcessor{
		pgPool:         pgPool,
		transactionsCh: make(chan *Transaction, jobsSize),

		shutdownContext: ctx,
		shutdownCancel:  cancel,
	}

	processor.wg.Add(workersCount)
	for range workersCount {
		go processor.processTransactions()
	}

	return processor
}

// Tell current workers to stop and wait for them to finish
//
// Multiple calls are safe to perform
func (p *transactionProcessor) Close(ctx context.Context) error {
	p.shutdownCancel()
	workersDone := make(chan struct{})
	go func() {
		p.wg.Wait()
		workersDone <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-workersDone:
	}

	return nil
}

// Insert adds a transaction to the processing queue
func (p *transactionProcessor) Insert(t *Transaction) {
	p.transactionsCh <- t
}

// TODO: Consider batching transactions?
func (p *transactionProcessor) processTransactions() {
	defer p.wg.Done()

	for {
		select {
		case t := <-p.transactionsCh:
			// TODO: Acquire a dedicated connection for each worker?
			if _, err := p.pgPool.Exec(
				context.Background(),
				"INSERT INTO transactions (user_id, transaction_type, amount, timestamp) VALUES ($1, $2, $3, $4);",
				t.UserID,
				t.TransactionType,
				t.Amount,
				t.Timestamp,
			); err != nil {
				log.Printf("Failed to save transaction in the database: %v", err)
			}
			fmt.Printf("%+v\n", t)
		case <-p.shutdownContext.Done():
			return
		}
	}
}
