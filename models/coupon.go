package models

import (
	"gorm.io/gorm"
)

type Coupons struct {
	gorm.Model
	Code              string  `gorm:"not null;unique"`
	Discount          float64 `gorm:"not null"`
	ExpiryDate        int     `gorm:"not null"` // In days
	MinAmount         float64 `gorm:"not null"`
	MaxDiscountAmount float64 `gorm:"not null"`
}
