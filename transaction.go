package main

import "github.com/jinzhu/gorm"

// Transaction model
type Transaction struct {
	gorm.Model
	UserID  uint
	GameID  uint
	Amount  int  `gorm:"not null"`
	IsBuyin bool `gorm:"not null;default:0"`
}
