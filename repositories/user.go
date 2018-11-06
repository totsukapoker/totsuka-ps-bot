package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/totsukapoker/totsuka-ps-bot/models"
)

// UserRepository has db property.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates new UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FirstOrCreate returns user by UserID if exists otherwise create user.
func (u *UserRepository) FirstOrCreate(UserID, DisplayName, PictureURL, StatusMessage string) (user models.User) {
	u.db.Where(models.User{UserID: UserID}).Assign(models.User{DisplayName: DisplayName, PictureURL: PictureURL, StatusMessage: StatusMessage}).FirstOrCreate(&user)
	return
}

// FindByIDs returns users by IDs.
func (u *UserRepository) FindByIDs(ids []uint) (users []models.User) {
	u.db.Where("ID in (?)", ids).Find(&users)
	return
}

// SetMyName changes user's MyName to name.
func (u *UserRepository) SetMyName(user *models.User, name string) {
	if user.MyName == name {
		return
	}
	user.MyName = name
	u.db.Save(&user)
}
