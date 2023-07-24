package models

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	UserID         string `gorm:"user_id"` // foregin key
	Name           string `gorm:"name,not null"`
	Mobile         string `gorm:"mobile,not null"`
	Pincode        string `gorm:"pincode,not null"`
	Locality       string `gorm:"locality,not null"`
	Address        string `gorm:"address,not null"`
	City           string `gorm:"city,not null"`
	State          string `gorm:"state,not null"`
	Landmark       string `gorm:"landmark"`
	AlternatePhone string `gorm:"alternate_phone"`
	IsDefault      bool   `gorm:"is_default,not null"`
}
