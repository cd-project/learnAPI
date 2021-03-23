package model

type User struct {
	ID       int    `gorm:"primaryKey"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}
