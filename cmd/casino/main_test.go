package main_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	casino "github.com/ncmprbll/altenar-assignment/cmd/casino"
)

var (
	app       *casino.App
	mock      sqlmock.Sqlmock
	mockKafka *mockKafkaConsumerImpl
)

type mockKafkaConsumerImpl struct {
	messages chan []byte
}

func (m *mockKafkaConsumerImpl) ReadMessage(timeout time.Duration) (*kafka.Message, error) {
	for {
		select {
		case <-time.After(timeout):
			return nil, kafka.NewError(kafka.ErrTimedOut, kafka.ErrTimedOut.String(), false)
		case msg := <-m.messages:
			return &kafka.Message{
				Value: msg,
			}, nil
		}
	}
}

func (m *mockKafkaConsumerImpl) SendMessage(msg []byte) {
	m.messages <- msg
}

func TestMain(m *testing.M) {
	re := regexp.MustCompile(`[\n\t ]`)
	db, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherFunc(func(expected, actual string) error {
		if expected == "skip" {
			return nil
		}
		if re.ReplaceAllString(actual, "") == re.ReplaceAllString(expected, "") {
			return nil
		}
		return errors.New("sql mismatch")
	})))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	mock = sqlMock
	mockKafka = &mockKafkaConsumerImpl{
		messages: make(chan []byte),
	}

	app = casino.NewApp(db, mockKafka)

	m.Run() // code := m.Run()

	// https://github.com/golang/go/commit/2f54081adfc967836842c96619d241378400ece6
	// os.Exit(code)
}
