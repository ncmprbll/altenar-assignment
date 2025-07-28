package main

import (
	"database/sql"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type app struct {
	db    *sql.DB
	kafka *kafka.Consumer
}
