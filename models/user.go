package models

import "github.com/jinzhu/gorm"

// User model
type User struct {
	gorm.Model
	UserID        string `gorm:"unique;not null"`
	DisplayName   string `gorm:"not null"`
	PictureURL    string
	StatusMessage string
	MyName        string
	Transactions  []Transaction
}

// Name - Get name of user on this service
func (user *User) Name() string {
	if user.MyName != "" {
		return user.MyName
	}
	return user.DisplayName
}
