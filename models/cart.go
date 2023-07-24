package models

import "gorm.io/gorm"

type UserCart struct {
	gorm.Model
	UserID       string `gorm:"user_id,not null"`
	ProductID    uint   `gorm:"product_id,not null"`
	ProductCount int    `gorm:"product_count,default:1,not null"`
}
