package main

import (
	"gateway/app/gateway"
	"gateway/app/mainPage"
	"gateway/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	gatewayRepository := gateway.NewHTTPRepository()
	gatewayService := gateway.NewService(gatewayRepository)
	gatewayController := gateway.NewController(gatewayService)

	r := gin.Default()
	r.LoadHTMLGlob("public/views/*")
	api := r.Group("/api")

	//Main Page
	r.GET("/", func(c *gin.Context) { mainPage.RenderMainPage(c) })

	//API routes
	api.Use(middlewares.BearerTokenMiddleware())
	{
		api.GET("/:service/*route", func(c *gin.Context) { gatewayController.Get(c) })
		api.POST("/:service/*route", func(c *gin.Context) { gatewayController.Post(c) })
		api.PUT("/:service/*route", func(c *gin.Context) { gatewayController.WithBody(c) })
		api.PATCH("/:service/*route", func(c *gin.Context) { gatewayController.WithBody(c) })
		api.DELETE("/:service/*route", func(c *gin.Context) { gatewayController.WithBody(c) })
	}

	err := r.Run(":9000")
	if err != nil {
		return
	}
}
