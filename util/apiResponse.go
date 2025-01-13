package util

import "github.com/gin-gonic/gin"

func SendError(c *gin.Context, statusCode int, errorMsg string) {
	c.JSON(statusCode, gin.H{"error": errorMsg})
	c.Abort()
}
