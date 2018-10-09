package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Game model
type Game struct {
	gorm.Model
	Name         string `gorm:"not null"`
	StartedAt    *time.Time
	EndedAt      *time.Time
	Transactions []Transaction
}
