package domain

import (
	"time"
)

// Loan represents the loans table in the database
type Loan struct {
	ID           int64     `gorm:"column:id;primaryKey" json:"id"`
	CustomerID   int64     `gorm:"column:customer_id;not null" json:"customer_id"`
	Principal    int64     `gorm:"column:principal;not null" json:"principal"`         // NUMERIC(15,0) → float64
	InterestRate float64   `gorm:"column:interest_rate;not null" json:"interest_rate"` // NUMERIC(5,2) → float64
	TotalAmount  int64     `gorm:"column:total_amount;not null" json:"total_amount"`   // NUMERIC(15,0) → float64
	TermWeeks    int       `gorm:"column:term_weeks;not null" json:"term_weeks"`
	StartDate    time.Time `gorm:"column:start_date;not null" json:"start_date"`
	Status       string    `gorm:"column:status;size:20;default:ONGOING" json:"status"`

	CreatedAt time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	CreatedBy string     `gorm:"column:created_by;size:255;not null" json:"created_by"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	UpdatedBy *string    `gorm:"column:updated_by;size:255" json:"updated_by"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	DeletedBy *string    `gorm:"column:deleted_by;size:255" json:"deleted_by"`
}

func (r *Loan) TableName() string {
	return "loans"
}

// Loan status enumeration
const (
	LoanStatus_Ongoing string = "ONGOING"
	LoanStatus_Paid    string = "PAID"
)

type FindOneLoanParam struct {
	CustomerID int64
	LoanID     int64
	Status     string
}
