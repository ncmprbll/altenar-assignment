package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	environment, err := godotenv.Read(".env")
	if err != nil {
		panic(err)
	}

	pgPool, err := pgxpool.New(ctx, environment["POSTGRES_DSN"])
	if err != nil {
		panic(err)
	}
	defer pgPool.Close()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": environment["KAFKA_BOOTSTRAP_SERVERS"],
		"group.id":          environment["KAFKA_CONSUMER_GROUP"],
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	if err := consumer.SubscribeTopics([]string{environment["KAFKA_TOPIC"]}, nil); err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	eventsSize := 100
	events := make(chan Transaction, eventsSize)

	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				log.Println("Consumer goroutine finished")
				return
			default:
				msg, err := consumer.ReadMessage(time.Second)
				if err == nil {
					log.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
					t := Transaction{}
					if err := json.Unmarshal(msg.Value, &t); err != nil {
						log.Printf("Failed to unmarshal message from the queue")
						continue
					}
					if !t.HasAllFields() {
						log.Printf("Bad message: not all fields are present")
						continue
					}
					events <- t
				} else if !err.(kafka.Error).IsTimeout() {
					// The client will automatically try to recover from all errors.
					// Timeout is not considered an error because it is raised by
					// ReadMessage in absence of messages.
					log.Printf("Consumer error: %v (%v)\n", err, msg)
				}
			}
		}
	}()

	go func() {
		defer wg.Done()

		for {
			select {
			case transaction := <-events:
				if _, err := pgPool.Exec(
					ctx,
					"INSERT INTO transactions (user_id, transaction_type, amount, timestamp) VALUES ($1, $2, $3, $4)",
					transaction.UserID,
					transaction.TransactionType,
					transaction.Amount,
					transaction.Timestamp,
				); err != nil {
					log.Println("Failed to save transaction in the database")
				}
			case <-ctx.Done():
				log.Println("Processing goroutine finished")
				return
			}
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Stopping goroutines...")
	// Tell everyone to finish their work and wrap up
	cancel()
	wg.Wait()
	log.Println("Consumers and processors have been stopped")

	_, gracefulCancel := context.WithTimeout(context.Background(), time.Minute)
	defer gracefulCancel()

	log.Println("Kafka consumer shutdown")
	consumer.Close()

	log.Println("Postgres pool shutdown")
	pgPool.Close()
}
