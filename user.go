package main

import "github.com/jinzhu/gorm"

// User model
type User struct {
	gorm.Model
	UserID        string `gorm:"unique;not null"`
	DisplayName   string `gorm:"not null"`
	PictureURL    string
	StatusMessage string
	Transactions  []Transaction
}
