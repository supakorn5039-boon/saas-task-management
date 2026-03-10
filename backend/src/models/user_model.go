package models

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
