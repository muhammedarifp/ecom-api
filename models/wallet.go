package models

import "gorm.io/gorm"

type Wallets struct {
	gorm.Model
	UserID  uint    `gorm:"user_id"`
	Balance float64 `gorm:"balance"`
}
