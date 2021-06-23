package model

import "gorm.io/gorm"

// 用户的基础信息
type UserBase struct {
	gorm.Model
	Username string `gorm:"uniqueIndex"`
	Password string
	Salt     string
}
