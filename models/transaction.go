package models

import "time"

// Transaction model
type Transaction struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	User      User
	UserID    uint
	Game      Game
	GameID    uint
	Amount    int  `gorm:"not null"`
	IsBuyin   bool `gorm:"not null;default:0"`
}
