package seeder

import (
	"log"

	"github.com/supakorn5039-boon/saas-task-backend/src/database/model"
	"github.com/supakorn5039-boon/saas-task-backend/src/security"
	"gorm.io/gorm"
)

func UserSeeder(db *gorm.DB) {
	hashedPassword, err := security.HashPassword("password123")
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	usersToSeed := []model.User{
		{
			Email:    "admin@example.com",
			Password: hashedPassword,
			Status:   1,
			Role:     "admin",
		},
		{
			Email:    "user@example.com",
			Password: hashedPassword,
			Status:   1,
			Role:     "user",
		},
	}

	for _, user := range usersToSeed {
		var count int64
		db.Model(&model.User{}).Where("email = ?", user.Email).Count(&count)

		if count == 0 {
			if err := db.Create(&user).Error; err != nil {
				log.Printf("failed to seed user %s: %v", user.Email, err)
			} else {
				log.Printf("User seeded: %s / password123", user.Email)
			}
		} else {
			log.Printf("User %s already exists, skipping.", user.Email)
		}
	}
}
