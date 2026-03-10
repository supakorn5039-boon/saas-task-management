package database

import (
	"fmt"
	"log"

	"github.com/supakorn5039-boon/saas-task-backend/src/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func Init(config *config.DatabaseConfig) {
	host := config.Host
	port := config.Port
	password := config.Password
	dbName := config.Database
	user := config.User

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Bangkok", host, port, user, password, dbName)

	var err error

	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
}
