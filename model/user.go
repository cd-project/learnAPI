package model

import "github.com/dgrijalva/jwt-go"

type User struct {
	ID                 int    `json:"id"       gorm:"primaryKey"`
	Username           string `json:"username" gorm:"column:username"`
	Password           string `json:"password" gorm:"column:password"`
	Role               string `json:"role" gorm:"column:role"`
	jwt.StandardClaims `gorm:"-"`
}

type UserRepository interface {
	GetAll() ([]UserResponse, error)
}
