package cmd

import (
	"gateway/app/jwt"
	"gateway/app/redis"
	"gateway/database/seeder"
	"gateway/util"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB    *gorm.DB
	Redis redis.Service
	Jwt   jwt.Service
}

func InitDependencies() *Dependencies {
	db, err := InitDB()
	if err != nil {
		util.SaveErrToFile(err)
	}
	redisService, err := InitRedis()
	if err != nil {
		util.SaveErrToDB(err, db)
	}
	jwtService, err := InitJwt()
	if err != nil {
		util.SaveErrToDB(err, db)
	}
	err = seeder.Seed(db)
	if err != nil {
		util.SaveErrToDB(err, db)
	}

	return &Dependencies{
		DB:    db,
		Redis: redisService,
		Jwt:   jwtService,
	}
}
