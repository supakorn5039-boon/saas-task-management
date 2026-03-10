package main

import (
	"fmt"
	"log"

	"github.com/supakorn5039-boon/saas-task-backend/src/config"
	"github.com/supakorn5039-boon/saas-task-backend/src/database/migrations"
	"github.com/supakorn5039-boon/saas-task-backend/src/database/seeders"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	appConfig := config.NewAppConfig()

	if err := appConfig.Load("config.ini"); err != nil {
		log.Fatal(err)
	}

	host := appConfig.Config.Database.Host
	port := appConfig.Config.Database.Port
	user := appConfig.Config.Database.User
	password := appConfig.Config.Database.Password
	dbname := appConfig.Config.Database.Database

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Bangkok",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	// Run migration
	migrations.Migrate(db)

	// Run seed
	seeders.Seed(db)
}
