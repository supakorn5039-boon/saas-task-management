package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string `gorm:"not null;unique"`
	Description string `gorm:"null"`
}
