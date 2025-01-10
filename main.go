package main

import (
	"gateway/app/gateway"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type GetStruct struct {
	Service string `uri:"service" binding:"required"`
	Route   string `uri:"route"`
}

func main() {
	godotenv.Load(".env")

	repository := gateway.NewHTTPRepository()
	service := gateway.NewService(repository)
	controller := gateway.NewController(service)

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	r.GET("/:service/*route", func(c *gin.Context) {
		controller.Get(c)
	})

	r.POST("/:service/*route", func(c *gin.Context) {

	})

	err := r.Run(":9000")
	if err != nil {
		return
	}
}
