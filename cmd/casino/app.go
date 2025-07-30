package main

import (
	"database/sql"
)

type App struct {
	DB    *sql.DB
	Kafka KafkaConsumerWrapper
}

func NewApp(db *sql.DB, kafka KafkaConsumerWrapper) *App {
	return &App{
		DB:    db,
		Kafka: kafka,
	}
}
