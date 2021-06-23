package model

import "gossodemo/internal/pkg/database"

func Init() {
	database.DB.AutoMigrate(&Product{})
	database.DB.AutoMigrate(&UserBase{})
}
