package util

import (
	"encoding/json"
	"gateway/database/model"
	"gorm.io/gorm"
	"os"
	"runtime"
	"runtime/debug"
)

func SaveErrToDB(err error, db *gorm.DB) {
	log := errToLog(err)

	db.Create(&log)
	return
}

func SaveErrToFile(err error) {
	log := errToLog(err)

	jsonLog, _ := json.Marshal(log)

	file, _ := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	defer file.Close()

	file.Write([]byte("***______________***\n"))

	file.Write(jsonLog)

	file.Write([]byte("\n***______________***\n"))
}

func errToLog(err error) model.SystemLog {
	log := model.SystemLog{
		Message: err.Error(),
		Stack:   string(debug.Stack()),
	}

	if _, file, line, _ := runtime.Caller(1); file != "" {
		log.File = file
		log.Line = line
	}

	return log
}
