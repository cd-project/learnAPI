package model

type Profile struct {
	ID     int    `gorm:"primaryKey"`
	Name   string `gorm:"column:name"`
	Age    int    `gorm:"column:age"`
	Gender string `gorm:"column:gender"`
	UID    int    `gorm<:"column:uid;unique"`
}
