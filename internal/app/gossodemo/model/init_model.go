package model

import (
    "gossodemo/internal/pkg/database"
    "log"
)

func Init() {
    db := database.DB
    errorList := make([]error, 0)
    errorList = appendError(errorList, db.AutoMigrate(&Product{}))
    errorList = appendError(errorList, db.AutoMigrate(&UserBase{}))

    if len(errorList) > 0 {
        for _, err := range errorList {
            log.Println(err)
        }

        panic("Init database fialed!")
    }
}

func appendError(errorList []error, err error) []error {
    if err != nil {
        return append(errorList, err)
    }
    
    return errorList
}
