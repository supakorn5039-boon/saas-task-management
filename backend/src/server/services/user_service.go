package services

import (
	"github.com/supakorn5039-boon/saas-task-backend/src/database"
	"github.com/supakorn5039-boon/saas-task-backend/src/models"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService() *UserService {
	return &UserService{database.Db}
}

func (s *UserService) GetUserById(id uint) (*models.UserDto, error) {
	var user *models.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return user.ToDto(), nil
}
