package domain

import (
	"time"
)

// Customer represents the customer table in the database
type Customer struct {
	ID          int64   `gorm:"column:id;primaryKey" json:"id"`
	Name        string  `gorm:"column:name;size:255;not null" json:"name"`
	PhoneNumber *string `gorm:"column:phone_number;size:50;not null" json:"phone_number"`

	CreatedAt time.Time  `gorm:"column:created_at;not null" json:"created_at"`
	CreatedBy string     `gorm:"column:created_by;size:255;not null" json:"created_by"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	UpdatedBy *string    `gorm:"column:updated_by;size:255" json:"updated_by"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	DeletedBy *string    `gorm:"column:deleted_by;size:255" json:"deleted_by"`
}

func (r *Customer) TableName() string {
	return "customers"
}

type FindAllProductFilterParam struct {
	ProductIDs  []int64
	Name        *string
	Description *string
}

type GetProductListParam struct {
	ProductIDs []int64
	Name       *string
}
