package seeder

import (
	"github.com/supakorn5039-boon/saas-task-backend/src/database/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func init() {
	Register(&seedRoles{})
}

type seedRoles struct{}

func (s *seedRoles) ID() string          { return "00_roles" }
func (s *seedRoles) Description() string { return "System roles: admin, manager, user" }

func (s *seedRoles) Seed(db *gorm.DB) error {
	roles := []model.Role{
		{Name: "admin", Description: "Full system access"},
		{Name: "manager", Description: "Can manage tasks and projects"},
		{Name: "user", Description: "Regular user access"},
	}

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoNothing: true,
	}).Create(&roles).Error
}
