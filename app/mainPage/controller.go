package mainPage

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime"
)

func RenderMainPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"GoVersion":  runtime.Version(),
		"GinVersion": gin.Version,
	})
	return
}
