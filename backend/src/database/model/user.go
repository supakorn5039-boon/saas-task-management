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
	Id        uint   `json:"id"`
	Email     string `json:"email"`
	Status    int    `json:"status"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
}

// Allowed roles — keep in sync with the seeder + frontend nav-config.
func IsValidRole(role string) bool {
	switch role {
	case "admin", "manager", "user":
		return true
	}
	return false
}

type UserListMeta struct {
	Page    int   `json:"page"`
	PerPage int   `json:"perPage"`
	Total   int64 `json:"total"`
}

type UserListResponse struct {
	Data []*UserDto   `json:"data"`
	Meta UserListMeta `json:"meta"`
}

func (u *User) ToDto() *UserDto {
	return &UserDto{
		Id:        u.ID,
		Email:     u.Email,
		Status:    u.Status,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
