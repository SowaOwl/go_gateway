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

	deps := cmd.InitDependencies()
	gateway := cmd.InitControllers()
	authMiddleware, logMiddleware := cmd.InitMiddlewares(deps)

	r := gin.Default()
	r.LoadHTMLGlob("public/views/*")
	api := r.Group("/api")

	//Main Page
	r.GET("/", func(c *gin.Context) { mainPage.RenderMainPage(c) })

	//API routes
	api.Use(authMiddleware.BearerTokenMiddleware(), logMiddleware.LogRequestMiddleware())
	{
		api.GET("/:service/*route", func(c *gin.Context) { gateway.Get(c) })
		api.POST("/:service/*route", func(c *gin.Context) { gateway.Post(c) })
		api.PUT("/:service/*route", func(c *gin.Context) { gateway.WithBody(c) })
		api.PATCH("/:service/*route", func(c *gin.Context) { gateway.WithBody(c) })
		api.DELETE("/:service/*route", func(c *gin.Context) { gateway.WithBody(c) })
	}

	err = r.Run(":9000")
	if err != nil {
		util.SaveErrToDB(err, deps.DB)
	}
}
