package service

import (
	"errors"
	"fmt"

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
		return nil, errors.New("invalid email or password")
	}

	if !security.VerifyPassword(user.Password, password) {
		return nil, errors.New("invalid email or password")
	}

	return user.ToDto(), nil
}

func (s *AuthService) Register(email, password string) (*model.UserDto, error) {
	var existing model.User
	if err := s.db.First(&existing, "email = ?", email).Error; err == nil {
		return nil, fmt.Errorf("email already exists")
	}

	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	newUser := model.User{
		Email:    email,
		Password: hashedPassword,
		Status:   1,
		Role:     "user",
	}

	if err := s.db.Create(&newUser).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	return newUser.ToDto(), nil
}
