package main

import (
	"gateway/app/gateway"
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

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"GoVersion":  runtime.Version(),
			"GinVersion": gin.Version,
		})
	})

	api.Use()
	{
		api.GET("/:service/*route", func(c *gin.Context) { controller.Get(c) })
		api.POST("/:service/*route", func(c *gin.Context) {})
	}

	err := r.Run(":9000")
	if err != nil {
		return
	}
}
