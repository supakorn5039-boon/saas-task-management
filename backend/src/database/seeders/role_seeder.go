package seeders

import (
	"log"

	"github.com/supakorn5039-boon/saas-task-backend/src/models"
	"gorm.io/gorm"
)

func RoleSeeder(db *gorm.DB) {
	rolesToSeed := []models.Role{
		{
			Name:        "admin",
			Description: "Full system access",
		},
		{
			Name:        "manager",
			Description: "Can manage tasks and projects",
		},
		{
			Name:        "user",
			Description: "Regular user access",
		},
	}

	for _, role := range rolesToSeed {
		var count int64
		db.Model(&models.Role{}).Where("name = ?", role.Name).Count(&count)

		if count == 0 {
			if err := db.Create(&role).Error; err != nil {
				log.Printf("failed to seed role %s: %v", role.Name, err)
			} else {
				log.Printf("Role seeded: %s", role.Name)
			}
		} else {
			log.Printf("Role %s already exists, skipping.", role.Name)
		}
	}
}
