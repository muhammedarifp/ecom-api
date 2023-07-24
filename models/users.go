package models

import "gorm.io/gorm"

type Users struct {
	gorm.Model
	FirstName  string `gorm:"first_name"`
	LastName   string `gorm:"last_name"`
	Email      string `gorm:"email"`
	Mobile     string `gorm:"mobile"`
	Password   string `gorm:"password"`
	Isadmin    bool   `gorm:"is_admin"`
	Status     bool   `gorm:"status"`
	Isverified bool   `gorm:"is_verified"` // 0 is unverified 1 is verified
}
