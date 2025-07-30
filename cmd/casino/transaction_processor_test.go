package main_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	casino "github.com/ncmprbll/altenar-assignment/cmd/casino"
)

var processorTests = []struct {
	jobsSize              int
	workersCount          int
	shouldErrorOnCreation bool
}{
	{
		jobsSize:              1,
		workersCount:          1,
		shouldErrorOnCreation: false,
	},
	{
		jobsSize:              0,
		workersCount:          0,
		shouldErrorOnCreation: true,
	},
	{
		jobsSize:              -1,
		workersCount:          0,
		shouldErrorOnCreation: true,
	},
	{
		jobsSize:              0,
		workersCount:          -1,
		shouldErrorOnCreation: true,
	},
	{
		jobsSize:              -1,
		workersCount:          -1,
		shouldErrorOnCreation: true,
	},
	{
		jobsSize:              -10,
		workersCount:          -10,
		shouldErrorOnCreation: true,
	},
	{
		jobsSize:              0,
		workersCount:          1,
		shouldErrorOnCreation: false,
	},
	{
		jobsSize:              2,
		workersCount:          2,
		shouldErrorOnCreation: false,
	},
	{
		jobsSize:              4,
		workersCount:          8,
		shouldErrorOnCreation: false,
	},
}

func TestNewTransactionProcessor(t *testing.T) {
	for i, tt := range processorTests {
		t.Run(fmt.Sprintf("TestCase_%02d", i), func(t *testing.T) {
			processor, err := casino.NewTransactionProcessor(app.DB, tt.jobsSize, tt.workersCount)
			if err != nil {
				if tt.shouldErrorOnCreation {
					return
				}
				t.Errorf("Failed to create transaction processor: %v", err)
			}

			ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
			defer cancel()

			if err := processor.Close(ctx); err != nil {
				t.Errorf("Failed to close transaction processor: %v", err)
			}
		})
	}
}
