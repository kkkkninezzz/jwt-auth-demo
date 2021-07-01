package database

import (
	"errors"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Connect(dsn string, tablePrefix string, singularTable bool) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix,   // table name prefix, table for `User` would be `t_users`
			SingularTable: singularTable, // use singular table name, table for `User` would be `user` with this option enabled
		},
	})
	if err != nil {
		panic(err)
	}

	DB = db
	log.Println("Connection Opened to Database")
}

// 是否为记录没有找到的err
func IsRecordNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
