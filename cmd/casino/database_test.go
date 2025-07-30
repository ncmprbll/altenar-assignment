package main_test

import (
	"fmt"
	"regexp"
	"testing"

	casino "github.com/ncmprbll/altenar-assignment/cmd/casino"
)

const (
	StatementNothing = `SELECT
							user_id,
							transaction_type,
							amount,
							timestamp
						FROM transactions;`

	StatementUserID = `SELECT
						   user_id,
						   transaction_type,
						   amount,
						   timestamp
					   FROM transactions
					   WHERE user_id = $1;`

	StatementTransactionType = `SELECT
									user_id,
									transaction_type,
									amount,
									timestamp
								FROM transactions
								WHERE transaction_type = $1;`

	StatementUserIDAndTransactionType = `SELECT
											 user_id,
											 transaction_type,
											 amount,
											 timestamp
										 FROM transactions
										 WHERE user_id = $1 AND transaction_type = $2;`
)

var queryTests = []struct {
	userID          string
	transactionType casino.TransactionType
	expectedStmt    string
	expectedArgs    []any
}{
	{
		userID:          "",
		transactionType: "",
		expectedStmt:    StatementNothing,
		expectedArgs:    []any{},
	},
	{
		userID:          "1",
		transactionType: "",
		expectedStmt:    StatementUserID,
		expectedArgs:    []any{"1"},
	},
	{
		userID:          "2",
		transactionType: "",
		expectedStmt:    StatementUserID,
		expectedArgs:    []any{"2"},
	},
	{
		userID:          "5",
		transactionType: "",
		expectedStmt:    StatementUserID,
		expectedArgs:    []any{"5"},
	},
	{
		userID:          "",
		transactionType: casino.TransactionTypeBet,
		expectedStmt:    StatementTransactionType,
		expectedArgs:    []any{casino.TransactionTypeBet},
	},
	{
		userID:          "",
		transactionType: casino.TransactionTypeWin,
		expectedStmt:    StatementTransactionType,
		expectedArgs:    []any{casino.TransactionTypeWin},
	},
	{
		userID:          "",
		transactionType: casino.TransactionTypeAll,
		expectedStmt:    StatementTransactionType,
		expectedArgs:    []any{casino.TransactionTypeAll},
	},
	{
		userID:          "1",
		transactionType: casino.TransactionTypeAll,
		expectedStmt:    StatementUserIDAndTransactionType,
		expectedArgs:    []any{"1", casino.TransactionTypeAll},
	},
	{
		userID:          "2",
		transactionType: casino.TransactionTypeAll,
		expectedStmt:    StatementUserIDAndTransactionType,
		expectedArgs:    []any{"2", casino.TransactionTypeAll},
	},
	{
		userID:          "3",
		transactionType: casino.TransactionTypeAll,
		expectedStmt:    StatementUserIDAndTransactionType,
		expectedArgs:    []any{"3", casino.TransactionTypeAll},
	},
	{
		userID:          "10",
		transactionType: casino.TransactionTypeAll,
		expectedStmt:    StatementUserIDAndTransactionType,
		expectedArgs:    []any{"10", casino.TransactionTypeAll},
	},
}

func TestBuildFindTransactionsQueryWithLogicalAnds(t *testing.T) {
	var re = regexp.MustCompile(`[\n\t ]`)
	for i, tt := range queryTests {
		t.Run(fmt.Sprintf("TestCase_%02d", i), func(t *testing.T) {
			stmt, args := app.BuildFindTransactionsQueryWithLogicalAnds(tt.userID, tt.transactionType)

			stmt = re.ReplaceAllString(stmt, "")
			expectedStmt := re.ReplaceAllString(tt.expectedStmt, "")

			if stmt != expectedStmt {
				t.Errorf("expected and resulting statement mismatch")
			}

			if len(args) != len(tt.expectedArgs) {
				t.Errorf("args len mismatch")
			}

			for k, v := range args {
				if v != tt.expectedArgs[k] {
					t.Errorf("args mismatch")
				}
			}
		})
	}
}
