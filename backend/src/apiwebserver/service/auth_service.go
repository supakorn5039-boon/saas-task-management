package service

import (
	"errors"

	"github.com/supakorn5039-boon/saas-task-backend/src/apperror"
	"github.com/supakorn5039-boon/saas-task-backend/src/database"
	"github.com/supakorn5039-boon/saas-task-backend/src/database/model"
	"github.com/supakorn5039-boon/saas-task-backend/src/security"
	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService() *AuthService {
	return &AuthService{database.DB}
}

func (s *AuthService) Login(email, password string) (*model.UserDto, error) {
	var user model.User
	if err := s.db.First(&user, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.Unauthorized("invalid email or password")
		}
		return nil, apperror.Wrap(err, 500, "login failed")
	}

	if !security.VerifyPassword(user.Password, password) {
		return nil, apperror.Unauthorized("invalid email or password")
	}

	return user.ToDto(), nil
}

func (s *AuthService) Register(email, password string) (*model.UserDto, error) {
	var existing model.User
	err := s.db.First(&existing, "email = ?", email).Error
	if err == nil {
		return nil, apperror.Conflict("email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.Wrap(err, 500, "register failed")
	}

	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		return nil, apperror.Wrap(err, 500, "register failed")
	}

	newUser := model.User{
		Email:    email,
		Password: hashedPassword,
		Status:   1,
		Role:     "user",
	}

	if err := s.db.Create(&newUser).Error; err != nil {
		return nil, apperror.Wrap(err, 500, "register failed")
	}

	return newUser.ToDto(), nil
}
