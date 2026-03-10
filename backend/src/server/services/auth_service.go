package services

import (
	"errors"
	"fmt"

	"github.com/supakorn5039-boon/saas-task-backend/src/database"
	"github.com/supakorn5039-boon/saas-task-backend/src/models"
	"github.com/supakorn5039-boon/saas-task-backend/src/security"
	"gorm.io/gorm"
)

type AuthenticationService struct {
	db *gorm.DB
}

func NewAuthenticationService() *AuthenticationService {
	return &AuthenticationService{database.Db}
}

func (s *AuthenticationService) Login(email, password string) (*models.UserDto, error) {
	var user models.User
	if err := s.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !security.VerifyPassword(user.Password, password) {
		return nil, errors.New("invalid email or password")
	}

	return user.ToDto(), nil
}

func (s *AuthenticationService) Register(email, password string) (*models.UserDto, error) {
	var existing models.User
	if err := s.db.First(&existing, "email = ?", email).Error; err == nil {
		return nil, fmt.Errorf("email already exists")
	}

	hashedPassword, err := security.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	newUser := models.User{
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
