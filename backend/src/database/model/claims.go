package model

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	jwt.RegisteredClaims
	Id    uint   `json:"id"`
	Role  string `json:"role"`
	Email string `json:"email,omitempty"`
}
