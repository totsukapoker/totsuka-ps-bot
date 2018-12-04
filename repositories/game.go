package repositories

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/totsukapoker/totsuka-ps-bot/models"
)

// GameRepository has db property.
type GameRepository struct {
	db *gorm.DB
}

// NewGameRepository creates new UserRepository.
func NewGameRepository(db *gorm.DB) *GameRepository {
	return &GameRepository{db: db}
}

// First returns game by id.
func (g *GameRepository) First(id uint) (game models.Game) {
	g.db.Preload("Transactions").Preload("Transactions.User").First(&game, id)
	return
}

// Current returns game which running right now.
func (g *GameRepository) Current() (game models.Game) {
	g.db.Where("? BETWEEN started_at AND ended_at", time.Now()).First(&game)
	return
}
