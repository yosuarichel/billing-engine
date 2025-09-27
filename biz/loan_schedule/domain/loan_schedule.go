package domain

import (
	"time"
)

// LoanSchedule represents the loan_schedule table in the database
type LoanSchedule struct {
	ID         int64     `gorm:"column:id;primaryKey" json:"id"`
	LoanID     int64     `gorm:"column:loan_id;not null" json:"loan_id"`
	WeekNumber int       `gorm:"column:week_number;not null" json:"week_number"`
	DueDate    time.Time `gorm:"column:due_date;not null" json:"due_date"`
	Amount     int64     `gorm:"column:amount;not null" json:"amount"` // NUMERIC(15, 0) → float64
	IsPaid     bool      `gorm:"column:is_paid;default:false" json:"is_paid"`

	CreatedAt time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	CreatedBy string     `gorm:"column:created_by;size:255;not null" json:"created_by"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	UpdatedBy *string    `gorm:"column:updated_by;size:255" json:"updated_by"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	DeletedBy *string    `gorm:"column:deleted_by;size:255" json:"deleted_by"`
}

func (*LoanSchedule) TableName() string {
	return "loan_schedule"
}

type FindAllScheduleParam struct {
	LoanScheduleID int64
	LoanID         int64
	WeekNumber     int
	IsPaid         *bool
}

type UpdateLoanSchedulesToPaidParam struct {
	LoanID     int64
	WeekNumber int
	IsPaid     *bool
}
