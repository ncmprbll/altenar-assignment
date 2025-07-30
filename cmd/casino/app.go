package main

import (
	"database/sql"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type App struct {
	db    *sql.DB
	kafka *kafka.Consumer
}

func NewApp(db *sql.DB, kafka *kafka.Consumer) *App {
	return &App{
		db:    db,
		kafka: kafka,
	}
}
