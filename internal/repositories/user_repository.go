package repositories

import (
	"github.com/thienhi/fusionstart/internal/dto"
	"github.com/thienhi/fusionstart/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user dto.UserRegisterDTO) error
	FindByEmail(email string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(user dto.UserRegisterDTO) error {
	newUser := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	}
	if err := r.db.Create(&newUser).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user *models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
