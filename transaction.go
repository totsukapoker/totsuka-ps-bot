package main

import "github.com/jinzhu/gorm"

// Transaction model
type Transaction struct {
	gorm.Model
	UserID uint
	amount int
}
