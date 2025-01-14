package main

import (
	"gateway/app/gateway"
	"gateway/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
	"runtime"
)

func main() {
	godotenv.Load(".env")

	repository := gateway.NewHTTPRepository()
	service := gateway.NewService(repository)
	controller := gateway.NewController(service)

	r := gin.Default()
	r.LoadHTMLGlob("public/views/*")
	api := r.Group("/api")

	//TODO убрать в свой контроллер
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"GoVersion":  runtime.Version(),
			"GinVersion": gin.Version,
		})
	})

	api.Use(middlewares.BearerTokenMiddleware())
	{
		api.GET("/:service/*route", func(c *gin.Context) { controller.Get(c) })
		api.POST("/:service/*route", func(c *gin.Context) { controller.Post(c) })
	}

	err := r.Run(":9000")
	if err != nil {
		return
	}
}
