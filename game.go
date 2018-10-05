package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Game model
type Game struct {
	gorm.Model
	Name         string     `gorm:"not null"`
	StartedAt    *time.Time `gorm:"not null"`
	EndedAt      *time.Time `gorm:"not null"`
	Transactions []Transaction
}
