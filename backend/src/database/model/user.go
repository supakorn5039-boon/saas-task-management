package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
	Status   int    `gorm:"default:0"`
	Role     string `gorm:"not null"`
}

type UserDto struct {
	Id     uint   `json:"id"`
	Email  string `json:"email"`
	Status int    `json:"status"`
	Role   string `json:"role"`
}

type CredentialDto struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (u *User) ToDto() *UserDto {
	return &UserDto{
		Id:     u.ID,
		Email:  u.Email,
		Status: u.Status,
		Role:   u.Role,
	}
}
