package domain

import (
	"time"
)

// Payment represents the payment table in the database
type Payment struct {
	ID          int64     `gorm:"column:id;primaryKey" json:"id"`
	LoanID      int64     `gorm:"column:loan_id;not null" json:"loan_id"`
	ScheduleID  int64     `gorm:"column:schedule_id;not null" json:"schedule_id"`
	PaymentDate time.Time `gorm:"column:payment_date;not null" json:"payment_date"`
	Amount      int64     `gorm:"column:amount;not null" json:"amount"`

	CreatedAt time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	CreatedBy string     `gorm:"column:created_by;size:255;not null" json:"created_by"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	UpdatedBy *string    `gorm:"column:updated_by;size:255" json:"updated_by"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	DeletedBy *string    `gorm:"column:deleted_by;size:255" json:"deleted_by"`
}

func (*Payment) TableName() string {
	return "payments"
}

type FindAllParam struct {
	PaymentID   int64
	LoanID      int64
	ScheduleID  int64
	PaymentDate string
}
