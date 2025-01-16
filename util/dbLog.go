package util

import (
	"gateway/database/model"
	"gorm.io/gorm"
	"runtime"
	"runtime/debug"
)

func SaveErrToDB(err error, db *gorm.DB) {
	log := model.SystemLog{
		Message: err.Error(),
		Stack:   string(debug.Stack()),
	}

	if _, file, line, _ := runtime.Caller(1); file != "" {
		log.File = file
		log.Line = line
	}

	db.Create(&log)
	return
}
