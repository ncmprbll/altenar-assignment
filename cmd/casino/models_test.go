package main_test

import (
	"fmt"
	"testing"
	"time"

	casino "github.com/ncmprbll/altenar-assignment/cmd/casino"
)

func makeshiftInt() *int {
	var a int
	return &a
}

func makeshiftType() *casino.TransactionType {
	var a casino.TransactionType
	return &a
}

func makeshiftTime() *time.Time {
	var a time.Time
	return &a
}

var hasAllFieldsTests = []struct {
	transaction casino.Transaction
	expected    bool
}{
	{
		casino.Transaction{},
		false,
	},

	{
		casino.Transaction{
			UserID: makeshiftInt(),
		},
		false,
	},
	{
		casino.Transaction{
			TransactionType: makeshiftType(),
		},
		false,
	},
	{
		casino.Transaction{
			Amount: makeshiftInt(),
		},
		false,
	},
	{
		casino.Transaction{
			Timestamp: makeshiftTime(),
		},
		false,
	},

	{
		casino.Transaction{
			UserID:          makeshiftInt(),
			TransactionType: makeshiftType(),
		},
		false,
	},
	{
		casino.Transaction{
			TransactionType: makeshiftType(),
			Amount:          makeshiftInt(),
		},
		false,
	},
	{
		casino.Transaction{
			Amount:    makeshiftInt(),
			Timestamp: makeshiftTime(),
		},
		false,
	},

	{
		casino.Transaction{
			UserID:          makeshiftInt(),
			TransactionType: makeshiftType(),
			Amount:          makeshiftInt(),
		},
		false,
	},
	{
		casino.Transaction{
			TransactionType: makeshiftType(),
			Amount:          makeshiftInt(),
			Timestamp:       makeshiftTime(),
		},
		false,
	},

	{
		casino.Transaction{
			UserID:          makeshiftInt(),
			TransactionType: makeshiftType(),
			Amount:          makeshiftInt(),
			Timestamp:       makeshiftTime(),
		},
		true,
	},
}

func TestMethodHasAllFields(t *testing.T) {
	for i, tt := range hasAllFieldsTests {
		t.Run(fmt.Sprintf("TestCase_%02d", i), func(t *testing.T) {
			got := tt.transaction.HasAllFields()
			if got != tt.expected {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}

var validTests = []struct {
	transaction casino.Transaction
	expected    bool
}{
	{
		casino.NewTransaction(0, "", 0, time.Now()),
		false,
	},
	{
		casino.NewTransaction(1, "", 0, time.Now()),
		false,
	},
	{
		casino.NewTransaction(0, "", 1, time.Now()),
		false,
	},
	{
		casino.NewTransaction(-1, "", 0, time.Now()),
		false,
	},
	{
		casino.NewTransaction(0, "", -1, time.Now()),
		false,
	},
	{
		casino.NewTransaction(-1, "", -1, time.Now()),
		false,
	},
	{
		casino.NewTransaction(-5, "", -5, time.Now()),
		false,
	},
	{
		casino.NewTransaction(-10, "", 5, time.Now()),
		false,
	},
	{
		casino.NewTransaction(10, "", -5, time.Now()),
		false,
	},
	{
		casino.NewTransaction(1, "", 1, time.Now()),
		true,
	},
	{
		casino.NewTransaction(5, "", 5, time.Now()),
		true,
	},
	{
		casino.NewTransaction(25, "", 25, time.Now()),
		true,
	},
	{
		casino.NewTransaction(25, "", 25, time.Now()),
		true,
	},
}

func TestMethodValid(t *testing.T) {
	for i, tt := range validTests {
		t.Run(fmt.Sprintf("TestCase_%02d", i), func(t *testing.T) {
			got := tt.transaction.Valid()
			if got != tt.expected {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}
