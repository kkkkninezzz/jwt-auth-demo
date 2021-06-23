package model

import "gossodemo/internal/pkg/database"

func Init() {
	db := database.DB
	db.AutoMigrate(&Product{})
	db.AutoMigrate(&UserBase{})
}
