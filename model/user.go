package model

import "github.com/dgrijalva/jwt-go"

const DefaultPassword = "0000"

type User struct {
	ID                 int    `json:"id"       gorm:"primaryKey"`
	Username           string `json:"username" gorm:"column:username"`
	Password           string `json:"password" gorm:"column:password"`
	Role               string `json:"role"     gorm:"column:role"`
	jwt.StandardClaims `gorm:"-"`
}

type UserResponse struct {
	ID       int    `json:"id"       gorm:"primaryKey"`
	Username string `json:"username" gorm:"column:username"`
	Role     string `json:"role"     gorm:"column:role"`
}

type UserPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserList struct {
	UserList []UserPayload `json:"userList"`
}
type UserRepository interface {
	GetAll() ([]UserResponse, error)
	GetByID(id int) (*User, error)
	GetByUsername(username string) (*User, error)
	Insert(newUser *User) (*User, error)
	ChangePassword(id int, newPwd string) error
	ChangeRole(id int, newRole string) (*User, error)
	DeleteUser(id int) error
	LoginTokenRequest(*User) (bool, error)
}
