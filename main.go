package main

import (
	"gateway/app/gateway"
	"gateway/app/mainPage"
	"gateway/cmd"
	"gateway/database/seeder"
	"gateway/middlewares"
	"gateway/util"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		util.SaveErrToFile(err)
	}

	db, _ := cmd.InitDB()
	redis, err := cmd.InitRedis()
	if err != nil {
		util.SaveErrToDB(err, db)
	}
	jwt, err := cmd.InitJwt()
	if err != nil {
		util.SaveErrToDB(err, db)
	}
	err = seeder.Seed(db)
	if err != nil {
		util.SaveErrToDB(err, db)
	}

	gatewayRepository := gateway.NewHTTPRepository()
	gatewayService := gateway.NewService(gatewayRepository)
	gatewayController := gateway.NewController(gatewayService)

	authMiddleware := middlewares.NewAuthMiddleware(db, jwt, redis)

	r := gin.Default()
	r.LoadHTMLGlob("public/views/*")
	api := r.Group("/api")

	//Main Page
	r.GET("/", func(c *gin.Context) { mainPage.RenderMainPage(c) })

	//API routes
	api.Use(authMiddleware.BearerTokenMiddleware(), middlewares.LogRequestMiddleware(db))
	{
		api.GET("/:service/*route", func(c *gin.Context) { gatewayController.Get(c) })
		api.POST("/:service/*route", func(c *gin.Context) { gatewayController.Post(c) })
		api.PUT("/:service/*route", func(c *gin.Context) { gatewayController.WithBody(c) })
		api.PATCH("/:service/*route", func(c *gin.Context) { gatewayController.WithBody(c) })
		api.DELETE("/:service/*route", func(c *gin.Context) { gatewayController.WithBody(c) })
	}

	err = r.Run(":9000")
	if err != nil {
		util.SaveErrToDB(err, db)
	}
}
