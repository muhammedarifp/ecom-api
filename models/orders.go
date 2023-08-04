package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID        uint    `gorm:"userid"`
	TottalAmount  float64 `gorm:"total_amount"`
	Status        string  `gorm:"status"`
	PayMethod     string  `gorm:"payment_method"`
	TransactionID uint    `gorm:"transaction_id"`
	IsSuccess     bool    `gorm:"is_success"`
	Items         []OrderItem
}

type OrderItem struct {
	gorm.Model
	OrderID   uint // foregin key
	ProductID uint // foregin key
	Quntity   uint
	Price     float64 // `gorm:"foreignkey:UserID"`
	Product   Product
}

// type
