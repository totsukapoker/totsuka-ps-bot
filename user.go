package main

import "github.com/jinzhu/gorm"

// User model
type User struct {
	gorm.Model
	UserID        string `gorm:"UNIQUE;NOT NULL"`
	DisplayName   string
	PictureURL    string
	StatusMessage string
	Transactions  []Transaction
}
