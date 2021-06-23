package model

import "gorm.io/gorm"

// 用户的基础信息
type UserBase struct {
	gorm.Model
	Username string `gorm:"VARCHAR(20) not null unique"`
	Password string `gorm:"VARCHAR(50) not null"`
	Salt     string `gorm:"VARCHAR(50) not null"`
}
