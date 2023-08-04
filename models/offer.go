package models

import "gorm.io/gorm"

type CatogoryOffer struct {
	gorm.Model
	CatogoryID           uint    `gorm:"not null"`
	Expiry               uint    `gorm:"catogory_id;not null"`
	MinTransactionAmount float64 `gorm:"not null"`
	Dicount              float64 `gorm:"not null"`
	MaxDiscount          float64 `gorm:"not null"`
	Catogory             Catogory
}
