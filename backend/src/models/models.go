package models

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Claims struct {
	jwt.RegisteredClaims
	Id uint `json:"id"`
}

type User struct {
	gorm.Model
	Email    string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	Status   int    `gorm:"default:0"`
	Role     string `gorm:"not null"`
}

type Role struct {
	gorm.Model
	Name        string `gorm:"not null;unique"`
	Description string `gorm:"null"`
}

type Task struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string `gorm:"null"`
	Completed   bool   `gorm:"default:false"`
	UserID      uint   `gorm:"not null"`
	User        User   `gorm:"foreignKey:UserID"`
}
