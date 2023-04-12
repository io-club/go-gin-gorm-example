package model

import "gorm.io/gorm"

// User 结构体表示用户信息
type User struct {
	gorm.Model
	Username string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
}

func (user User) TableName() string {
	return "user"
}

func GetUserById(id int) (User, error) {
	var user User
	result := DB.First(&user, id)
	return user, result.Error
}
