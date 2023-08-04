package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name       string  `gorm:"name"`
	Disc       string  `gorm:"disc"`
	Price      float64 `gorm:"price"`
	Quntity    uint    `gorm:"quntity"`
	CatogaryID uint    `gorm:"catogary_id"`
	IsActive   bool    `gorm:"is_active"`
	IsDeleted  bool    `gorm:"is_deleted"`
}

// Product Image table

type ProductImages struct {
	gorm.Model
	ProductID uint   `gorm:"product_id"` // foregin key
	ImageName string `gorm:"image_name"`
	IsDefault bool   `gorm:"is_default"`
}

// Catogary table struct

type Catogory struct {
	gorm.Model
	Name string `gorm:"name"`
	Disc string `gorm:"disc"` // this is a forein key
}
