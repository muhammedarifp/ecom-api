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
	MaxUsage          uint    `gorm:"not null"`
}

type CouponUsages struct {
	gorm.Model
	CouponID   uint   `gorm:"not null"`
	UserID     string `gorm:"not null"`
	UsageCount uint   `gorm:"not null"`
}
