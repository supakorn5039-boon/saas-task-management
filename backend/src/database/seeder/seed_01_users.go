package seeder

import (
	"fmt"

	"github.com/supakorn5039-boon/saas-task-backend/src/database/model"
	"github.com/supakorn5039-boon/saas-task-backend/src/security"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func init() {
	Register(&seedUsers{})
}

type seedUsers struct{}

func (s *seedUsers) ID() string          { return "01_users" }
func (s *seedUsers) Description() string { return "Default users (password: password123)" }

func (s *seedUsers) Seed(db *gorm.DB) error {
	hash, err := security.HashPassword("password123")
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	users := []model.User{
		{Email: "admin@example.com", Password: hash, Status: 1, Role: "admin"},
		{Email: "user@example.com", Password: hash, Status: 1, Role: "user"},
	}

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoNothing: true,
	}).Create(&users).Error
}
