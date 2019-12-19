package state

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// ExpenseTableName is the name of the expense table.
const ExpenseTableName = "expense"

// Expense represents the expense table.
type Member struct {
  UUID          string       		`db:"uuid"`         // text
	Description   string          `db:"description"` 	// text
	CreatedAt     time.Time       `db:"created_at"`   // date
	Amount   	  	float64 				`db:"amount"`   	  // float64
	Currency      string    			`db:"currency"`     // text
  Employee      Employee        `db:"employee"`    	// employee
}

// ExpenseField is an enum type for Expense fields.
type ExpenseField string

// ExpenseFields
const (
	ExpenseFieldUUID            ExpenseField = "uuid"
	ExpenseFieldDescription     ExpenseField = "description"
	ExpenseFieldCreatedAt 		ExpenseField = "created_at"
	ExpenseFieldAmount   		ExpenseField = "amount"
	ExpenseFieldCurrency      	ExpenseField = "currency"
	ExpenseFieldEmployee        ExpenseField = "employee"
)

func (f ExpenseField) String() string {
	return string(f)
}

// ExpenseListFilter enables ordering and paging Expense results.
type ExpenseListFilter struct {
	OrderBy ExpenseField
	Skip    *uint64
	Take    *uint64
}

// ExpenseFilter enables Expense result filtering.
type MemberFilter struct {
	List          	*ExpenseListFilter
	UUID            *StringFilter
	Description 	*StringFilter
	CreatedAt     	*TimeFilter
	Amount   		*Float64Filter
	Currency      	*StringFilter
	Employee        *EmployeeFilter
}
