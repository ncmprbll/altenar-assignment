package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// TransactionInserter is the interface that wraps the Insert method
//
// Insert places transaction t in a processing queue
type TransactionInserter interface {
	Insert(t *Transaction)
}

type transactionConsumer struct {
	consumer *kafka.Consumer
	inserter TransactionInserter

	shutdownContext context.Context
	shutdownCancel  context.CancelFunc
	wg              sync.WaitGroup
}

// Initialize a new transaction consumer
func NewTransactionConsumer(consumer *kafka.Consumer, inserter TransactionInserter) *transactionConsumer {
	ctx, cancel := context.WithCancel(context.Background())

	c := &transactionConsumer{
		consumer: consumer,
		inserter: inserter,

		shutdownContext: ctx,
		shutdownCancel:  cancel,
	}

	c.wg.Add(1)
	go c.consumeTransactions()

	return c
}

// Tell current workers to stop and wait for them to finish
//
// Multiple calls are safe to perform
func (c *transactionConsumer) Close(ctx context.Context) error {
	c.shutdownCancel()
	workersDone := make(chan struct{})
	go func() {
		c.wg.Wait()
		workersDone <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-workersDone:
	}

	return nil
}

func (c *transactionConsumer) consumeTransactions() {
	defer c.wg.Done()

	for {
		select {
		default:
			msg, err := c.consumer.ReadMessage(time.Second)
			if err == nil {
				t := &Transaction{}
				if err := json.Unmarshal(msg.Value, t); err != nil {
					log.Printf("Failed to unmarshal message from the queue")
					continue
				}

				if !t.Valid() {
					log.Printf("Bad message: not valid")
					continue
				}

				c.inserter.Insert(t)
			} else if !err.(kafka.Error).IsTimeout() {
				log.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		case <-c.shutdownContext.Done():
			return
		}
	}
}
