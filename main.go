package main

import (
	"gateway/app/mainPage"
	"gateway/cmd"
	"gateway/util"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		util.SaveErrToFile(err)
	}

	app := cmd.InitApp()

	gateway := app.Controllers.Gateway

	r := gin.Default()
	r.LoadHTMLGlob("public/views/*")
	api := r.Group("/api")

	//Main Page
	r.GET("/", func(c *gin.Context) { mainPage.RenderMainPage(c) })

	//API routes
	api.Use(app.Middlewares.Auth.Handle(), app.Middlewares.Log.Handle())
	{
		api.GET("/:service/*route", gateway.Get)
		api.POST("/:service/*route", gateway.Post)
		api.PUT("/:service/*route", gateway.WithBody)
		api.PATCH("/:service/*route", gateway.WithBody)
		api.DELETE("/:service/*route", gateway.WithBody)
	}

	err = r.Run(":9000")
	if err != nil {
		util.SaveErrToDB(err, app.Dependencies.DB)
	}
}
