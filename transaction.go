package main

import "github.com/jinzhu/gorm"

// Transaction model
type Transaction struct {
	gorm.Model
	UserID uint `gorm:"NOT NULL"`
	Amount int  `gorm:"NOT NULL"`
}
