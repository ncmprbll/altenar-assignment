package main_test

import (
	"context"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	casino "github.com/ncmprbll/altenar-assignment/cmd/casino"
)

var (
	app *casino.App
)

func TestMain(m *testing.M) {
	mainCtx, mainCancel := context.WithCancel(context.Background())
	defer mainCancel()

	environment, err := godotenv.Read("../../.env")
	if err != nil {
		panic(err)
	}

	pgPool, err := pgxpool.New(mainCtx, environment["POSTGRES_DSN"])
	if err != nil {
		panic(err)
	}
	defer pgPool.Close()

	db := stdlib.OpenDBFromPool(pgPool)
	defer db.Close()

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

	app = casino.NewApp(db, kafka)

	m.Run() // code := m.Run()

	// https://github.com/golang/go/commit/2f54081adfc967836842c96619d241378400ece6
	// os.Exit(code)
}
