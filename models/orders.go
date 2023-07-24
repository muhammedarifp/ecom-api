package models

import "gorm.io/gorm"

type Orders struct {
	gorm.Model
	UserID        string  `gorm:"userid"`
	TottalAmount  float64 `gorm:"total_amount"`
	Status        string  `gorm:"status"`
	PayMethod     string  `gorm:"payment_method"`
	TransactionID uint    `gorm:"transaction_id"`
	IsSuccess     bool    `gorm:"is_success"`
	// OrderItems    []OrdersItems `gorm:"foreignKey:ProductID"`
}

type OrdersItems struct {
	gorm.Model
	OrderID   uint    `gorm:"order_id"`   // foregin key
	ProductID uint    `gorm:"product_id"` // foregin key
	Quntity   uint    `gorm:"quntity"`
	Price     float64 `gorm:"price"`
	// Product   Products `gorm:"foreignKey:ProductID"`
}
