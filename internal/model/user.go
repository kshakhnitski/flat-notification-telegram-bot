package model

type User struct {
	ID        int64  `gorm:"primaryKey"`
	Username  string `gorm:"varchar(255);unique"`
	FirstName string `gorm:"varchar(255)"`
	ChatID    int64  `gorm:"int"`
}
