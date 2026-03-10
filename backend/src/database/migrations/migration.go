package migrations

import (
	"log"

	"github.com/supakorn5039-boon/saas-task-backend/src/models"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	log.Println("Migrating database...")
	if err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Task{},
	); err != nil {
		log.Fatalf("Failed to auto migrate: %v", err)
	}
	log.Println("Database migrate Successfully!")
}
