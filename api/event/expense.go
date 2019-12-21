package event

import uuid "github.com/satori/go.uuid"

// ExpenseUpdatedTopic notifies subscribers that a Expense has been updated.
const ExpenseUpdatedTopic Topic = "expense-updated"
const ExpenseInsertedTopic = "expense-inserted"

// ExpenseUpdated contains the data for the ExpenseUpdatedTopic.
type ExpenseUpdated struct {
	ID uuid.UUID `json:"id"`
}

// ExpenseInsertedTopic notifies subscribers a Expense has been inserted.

// ExpenseInserted contains the data for the ExpenseInsertedTopic.
type ExpenseInserted struct {
	ID string `json:"id"`
}
