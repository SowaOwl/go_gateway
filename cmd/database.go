package cmd

import (
	"gateway/app/database/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

func InitDB() (*gorm.DB, error) {
	dsn := os.Getenv("DB_USERNAME") + ":" + os.Getenv("DB_PASSWORD") +
		"@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") +
		")/" + os.Getenv("DB_DATABASE") + "?parseTime=true"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&models.ApiLog{},
		&models.WhiteListedEndpoint{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
