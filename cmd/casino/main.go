package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	mainCtx, mainCancel := context.WithCancel(context.Background())
	defer mainCancel()

	environment, err := godotenv.Read(".env")
	if err != nil {
		panic(err)
	}

	pgPool, err := pgxpool.New(mainCtx, environment["POSTGRES_DSN"])
	if err != nil {
		panic(err)
	}
	defer pgPool.Close()

	kafka, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": environment["KAFKA_BOOTSTRAP_SERVERS"],
		"group.id":          environment["KAFKA_CONSUMER_GROUP"],
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		panic(err)
	}
	defer kafka.Close()

	if err := kafka.SubscribeTopics([]string{environment["KAFKA_TOPIC"]}, nil); err != nil {
		panic(err)
	}

	processor := NewTransactionProcessor(pgPool, 8, 1)
	defer processor.Close(mainCtx)

	consumer := NewTransactionConsumer(kafka, processor)
	defer consumer.Close(mainCtx)

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	gracefulCtx, gracefulCancel := context.WithTimeout(context.Background(), time.Minute)
	defer gracefulCancel()

	// Firstly, we stop our consumers and workers...
	log.Println("Transaction consumer shutdown")
	if err := consumer.Close(gracefulCtx); err != nil {
		log.Printf("Failed to stop transaction consumer: %v", err)
	}

	log.Println("Transaction processor shutdown")
	if err := processor.Close(gracefulCtx); err != nil {
		log.Printf("Failed to stop transaction processor: %v", err)
	}

	// And then we tell everyone to stop doing what they are doing
	mainCancel()

	log.Println("Kafka consumer shutdown")
	kafka.Close()

	log.Println("Postgres pool shutdown")
	pgPool.Close()
}
