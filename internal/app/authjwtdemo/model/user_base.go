package model

import "gorm.io/gorm"

// 用户的基础信息
type UserBase struct {
    gorm.Model
    Username string `gorm:"type:string;not null;size:20;unique"`
    Password string `gorm:"type:string;not null;size:100;"`
    Salt     string `gorm:"type:string;not null;size:100;"`
}
