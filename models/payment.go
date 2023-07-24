package models

import "gorm.io/gorm"

type Transactions struct {
	gorm.Model
	OrderID  string  `json:"id"`
	Amount   float64 `json:"amount"`
	Attempts float64 `json:"attempts"`
	Currency string  `json:"currency"`
	Receipt  string  `json:"receipt"`
	Status   string  `json:"status"`
}
