package util

import (
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status bool        `json:"status"`
	Error  string      `json:"error"`
	Data   interface{} `json:"data"`
}

func SendError(c *gin.Context, statusCode int, errorMsg string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Status: false,
		Error:  errorMsg,
		Data:   data,
	})
	c.Abort()
}
