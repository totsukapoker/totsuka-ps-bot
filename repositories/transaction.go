package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/totsukapoker/totsuka-ps-bot/models"
)

// TransactionRepository has db property.
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates new UserRepository.
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// FindByGame returns all Transactions for game.
func (t *TransactionRepository) FindByGame(game models.Game) (transactions []models.Transaction) {
	t.db.Model(&game).Order("id desc").Related(&transactions)
	return
}

// LastBy returns last record of user on specify game.
func (t *TransactionRepository) LastBy(user models.User, game models.Game) (transaction models.Transaction) {
	t.db.Where("user_id = ? AND game_id = ?", user.ID, game.ID).Order("id desc").First(&transaction)
	return
}

// CurrentAmountBy returns current amount of user on specify game.
func (t *TransactionRepository) CurrentAmountBy(user models.User, game models.Game) int {
	type Result struct {
		Total int
	}
	var r Result
	t.db.Table("transactions").Select("IFNULL(SUM(amount), 0) AS total").Where("user_id = ? AND game_id = ?", user.ID, game.ID).Scan(&r)
	return r.Total
}

// CurrentAmountBuyinBy return current amount of buyin of user on specify game.
func (t *TransactionRepository) CurrentAmountBuyinBy(user models.User, game models.Game) int {
	type Result struct {
		Total int
	}
	var r Result
	t.db.Table("transactions").Select("IFNULL(SUM(amount), 0) AS total").Where("user_id = ? AND game_id = ? AND is_buyin = ?", user.ID, game.ID, true).Scan(&r)
	return r.Total
}

// Create make new transaction record and returns it.
func (t *TransactionRepository) Create(user models.User, game models.Game, amount int, isBuyin bool) (transaction models.Transaction) {
	transaction.UserID = user.ID
	transaction.GameID = game.ID
	transaction.Amount = amount
	transaction.IsBuyin = isBuyin
	t.db.Create(&transaction)
	return
}

// Delete deletes specify transaction record.
func (t *TransactionRepository) Delete(transaction *models.Transaction) {
	t.db.Delete(transaction)
}
